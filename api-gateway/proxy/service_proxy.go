package proxy

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type ServiceConfig struct {
	Name        string `mapstructure:"name"`
	Host        string `mapstructure:"host"`
	Port        int    `mapstructure:"port"`
	HealthCheck string `mapstructure:"health_check"`
}

type ServiceProxy struct {
	services map[string]*ServiceConfig
	proxies  map[string]*httputil.ReverseProxy
}

func NewServiceProxy(services map[string]ServiceConfig) *ServiceProxy {
	sp := &ServiceProxy{
		services: make(map[string]*ServiceConfig),
		proxies:  make(map[string]*httputil.ReverseProxy),
	}

	for name, config := range services {
		serviceConfig := config
		sp.services[name] = &serviceConfig

		target := fmt.Sprintf("http://%s:%d", config.Host, config.Port)
		targetURL, _ := url.Parse(target)
		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		// 优化错误处理
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			logrus.Errorf("🔴 Proxy error for service %s: %v", name, err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)

			// 根据请求方法返回不同的错误信息
			errorMsg := `{"code":502,"message":"Service temporarily unavailable","success":false}`
			if r.Method == http.MethodGet {
				errorMsg = `{"code":502,"message":"Service unavailable, please try again later","success":false,"data":null}`
			}
			w.Write([]byte(errorMsg))
		}

		// 优化响应处理
		proxy.ModifyResponse = func(resp *http.Response) error {
			// 设置统一的响应头
			resp.Header.Set("X-Service", name)
			resp.Header.Set("X-Gateway", "api-gateway")
			resp.Header.Set("X-Request-ID", resp.Header.Get("X-Request-ID"))
			resp.Header.Set("Access-Control-Allow-Origin", "*")
			resp.Header.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			resp.Header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// 记录响应信息
			logrus.Infof("📨 Response from %s: %d %s - %s %s",
				name, resp.StatusCode, http.StatusText(resp.StatusCode),
				resp.Request.Method, resp.Request.URL.Path)

			// 根据HTTP方法和状态码进行特殊处理
			switch resp.StatusCode {
			case http.StatusNotFound:
				logrus.Warnf("Service %s returned 404 for %s %s", name, resp.Request.Method, resp.Request.URL.Path)
				// 对于GET请求的404，可以返回更友好的错误信息
				if resp.Request.Method == http.MethodGet {
					resp.Header.Set("X-Cache-Hint", "no-cache")
				}

			case http.StatusInternalServerError:
				logrus.Errorf("Service %s returned 500 for %s %s", name, resp.Request.Method, resp.Request.URL.Path)
				// 对于POST/PUT请求的500错误，添加重试提示
				if resp.Request.Method == http.MethodPost || resp.Request.Method == http.MethodPut {
					resp.Header.Set("Retry-After", "30")
				}

			case http.StatusTooManyRequests:
				logrus.Warnf("Service %s rate limited: %s %s", name, resp.Request.Method, resp.Request.URL.Path)
				resp.Header.Set("Retry-After", "60")

			case http.StatusCreated:
				// 对于创建成功的请求，添加Location头（如果适用）
				if resp.Request.Method == http.MethodPost && resp.Header.Get("Location") == "" {
					if id := extractResourceID(resp.Request.URL.Path); id != "" {
						resp.Header.Set("Location", fmt.Sprintf("%s/%s", resp.Request.URL.Path, id))
					}
				}

			case http.StatusNoContent:
				// 对于204响应，确保没有响应体
				resp.ContentLength = 0
				resp.Body = http.NoBody
			}

			// 根据HTTP方法优化缓存头
			switch resp.Request.Method {
			case http.MethodGet:
				if resp.StatusCode == http.StatusOK {
					// 为GET请求设置缓存头
					resp.Header.Set("Cache-Control", "public, max-age=60")
					resp.Header.Set("ETag", generateETag(resp))
				}
			case http.MethodPost, http.MethodPut, http.MethodDelete:
				// 对于写操作，建议客户端不要缓存
				resp.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
				resp.Header.Set("Pragma", "no-cache")
				resp.Header.Set("Expires", "0")
			}

			return nil
		}

		sp.proxies[name] = proxy
	}

	return sp
}

func (sp *ServiceProxy) ProxyToService(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := fmt.Sprintf("%d", time.Now().UnixNano())

		logrus.Infof("🚀 Proxying %s %s to %s service",
			c.Request.Method, c.Request.URL.Path, serviceName)

		proxy, ok := sp.proxies[serviceName]
		if !ok {
			logrus.Errorf("Service not found: %s", serviceName)
			c.JSON(http.StatusNotFound, gin.H{
				"code":    404,
				"message": "Service not available",
				"success": false,
				"data":    nil,
			})
			return
		}

		// 设置请求头
		c.Request.Header.Set("X-Forwarded-For", c.ClientIP())
		c.Request.Header.Set("X-Gateway", "api-gateway")
		c.Request.Header.Set("X-Request-ID", requestID)
		c.Request.Header.Set("X-Real-IP", c.ClientIP())

		// 根据HTTP方法设置不同的头信息
		switch c.Request.Method {
		case http.MethodGet:
			c.Request.Header.Set("X-Cache", "true")
		case http.MethodPost, http.MethodPut:
			c.Request.Header.Set("X-Idempotency-Key", requestID)
		case http.MethodDelete:
			logrus.Infof("Delete operation requested for %s", c.Request.URL.Path)
		}

		proxy.ServeHTTP(c.Writer, c.Request)

		latency := time.Since(start)
		logrus.Infof("✅ %s %s -> %s completed in %v",
			c.Request.Method, c.Request.URL.Path, serviceName, latency)
	}
}

func (sp *ServiceProxy) HealthCheck(serviceName string) error {
	service, ok := sp.services[serviceName]
	if !ok {
		return fmt.Errorf("service %s not found", serviceName)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			ResponseHeaderTimeout: 3 * time.Second,
		},
	}

	url := fmt.Sprintf("http://%s:%d%s", service.Host, service.Port, service.HealthCheck)

	logrus.Infof("Health checking service %s at %s", serviceName, url)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Health check failed for %s: %v", serviceName, err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Errorf("Error closing response body: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		logrus.Errorf("Health check failed for %s: status %d", serviceName, resp.StatusCode)
		return fmt.Errorf("service %s health check failed: status %d", serviceName, resp.StatusCode)
	}

	logrus.Infof("Health check passed for %s", serviceName)
	return nil
}

func (sp *ServiceProxy) GetServiceStatus() map[string]bool {
	status := make(map[string]bool)
	for name := range sp.services {
		status[name] = sp.HealthCheck(name) == nil
	}
	return status
}

// 辅助函数：生成简单的ETag
func generateETag(resp *http.Response) string {
	return fmt.Sprintf("\"%d-%d\"", time.Now().Unix(), resp.ContentLength)
}

// 辅助函数：从URL路径中提取资源ID
func extractResourceID(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// 启动定时健康检查
func (sp *ServiceProxy) StartHealthChecks(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			for name := range sp.services {
				err := sp.HealthCheck(name)
				if err != nil {
					logrus.Warnf("Service %s health check failed: %v", name, err)
				} else {
					logrus.Debugf("Service %s is healthy", name)
				}
			}
		}
	}()
}

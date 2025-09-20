package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"reading-microservices/api-gateway/handlers"
	gatewayMiddleware "reading-microservices/api-gateway/middleware"
	"reading-microservices/api-gateway/proxy"
	"reading-microservices/shared/config"
)

type GatewayConfig struct {
	Server    config.ServerConfig            `mapstructure:"server"`
	Database  config.DatabaseConfig          `mapstructure:"database"`
	Redis     config.RedisConfig             `mapstructure:"redis"`
	JWT       config.JWTConfig               `mapstructure:"jwt"`
	Consul    config.ConsulConfig            `mapstructure:"consul"`
	Services  map[string]proxy.ServiceConfig `mapstructure:"services"`
	RateLimit struct {
		RequestsPerMinute int `mapstructure:"requests_per_minute"`
		Burst             int `mapstructure:"burst"`
	} `mapstructure:"rate_limit"`
}

func main() {
	// 设置更详细的日志
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	cfg, err := loadGatewayConfig()
	if err != nil {
		logrus.Fatalf("Failed to load config: %v", err)
	}

	logrus.Infof("Loaded configuration: Server=%s:%d", cfg.Server.Host, cfg.Server.Port)

	rdb, err := initRedis(cfg.Redis)
	if err != nil {
		logrus.Warnf("Redis init failed: %v", err)
		rdb = nil
	} else {
		logrus.Info("Redis connected successfully")
	}

	// 记录服务配置
	for name, service := range cfg.Services {
		logrus.Infof("Service %s: %s:%d", name, service.Host, service.Port)
	}

	serviceProxy := proxy.NewServiceProxy(cfg.Services)
	rateLimiter := gatewayMiddleware.NewRateLimiter(rdb, cfg.RateLimit.RequestsPerMinute, cfg.RateLimit.Burst)
	gatewayHandler := handlers.NewGatewayHandler(serviceProxy)
	router := setupRouter(gatewayHandler, rateLimiter, cfg.JWT.Secret)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logrus.Infof("API Gateway starting on %s", addr)

	if err := router.Run(addr); err != nil {
		logrus.Fatalf("Failed to start server: %v", err)
	}
}

func loadGatewayConfig() (*GatewayConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg GatewayConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func initRedis(cfg config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	return rdb, err
}

func setupRouter(handler *handlers.GatewayHandler, rl *gatewayMiddleware.RateLimiter, jwtSecret string) *gin.Engine {
	router := gin.Default()
	router.Use(gin.Recovery())
	router.GET("/health", handler.Health)
	router.GET("/status", handler.ServiceStatus)

	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		auth.Use(rl.IPLimit(30)) // IP限流30次/分钟
		{
			// 使用通配符处理所有 auth 路径
			auth.Any("/*path", handler.ProxyService("user_service"))
		}

		// ------------------------
		// 用户服务接口 - 必须登录
		// ------------------------
		user := v1.Group("/user")
		user.Use(gatewayMiddleware.AuthMiddleware(jwtSecret))
		user.Use(rl.UserLimit(5000))
		{
			user.Any("/*path", handler.ProxyService("user_service"))
		}
		v1.Any("/content/*path", gatewayMiddleware.OptionalAuth(jwtSecret), rl.UserLimit(1000), handler.ProxyService("content_service"))
		v1.Any("/reading/*path", gatewayMiddleware.AuthMiddleware(jwtSecret), rl.UserLimit(2000), handler.ProxyService("reading_service"))
		v1.Any("/payment/*path", gatewayMiddleware.AuthMiddleware(jwtSecret), rl.UserLimit(100), rl.IPLimit(50), handler.ProxyService("payment_service"))
		v1.Any("/download/*path", gatewayMiddleware.AuthMiddleware(jwtSecret), gatewayMiddleware.AntiLeechMiddleware(), rl.UserLimit(50), handler.ProxyService("download_service"))
		v1.Any("/notification/*path", gatewayMiddleware.AuthMiddleware(jwtSecret), rl.UserLimit(200), handler.ProxyService("notification_service"))
		v1.Any("/admin/*path", gatewayMiddleware.AuthMiddleware(jwtSecret), gatewayMiddleware.RoleMiddleware("admin"), rl.UserLimit(500), handler.ProxyService("admin_service"))
	}
	return router
}

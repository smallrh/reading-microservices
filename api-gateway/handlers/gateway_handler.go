package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reading-microservices/api-gateway/proxy"
)

type GatewayHandler struct {
	serviceProxy *proxy.ServiceProxy
}

func NewGatewayHandler(serviceProxy *proxy.ServiceProxy) *GatewayHandler {
	return &GatewayHandler{serviceProxy: serviceProxy}
}

func (h *GatewayHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "api-gateway", "version": "1.0.0"})
}

func (h *GatewayHandler) ServiceStatus(c *gin.Context) {
	status := h.serviceProxy.GetServiceStatus()
	allHealthy := true
	for _, healthy := range status {
		if !healthy {
			allHealthy = false
			break
		}
	}
	httpStatus := http.StatusOK
	if !allHealthy {
		httpStatus = http.StatusServiceUnavailable
	}
	c.JSON(httpStatus, gin.H{"status": "ok", "services": status, "healthy": allHealthy})
}

func (h *GatewayHandler) ProxyService(service string) gin.HandlerFunc {
	return h.serviceProxy.ProxyToService(service)
}

package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func AntiLeechMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		referer := c.Request.Header.Get("Referer")
		if referer == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"code": 403, "message": "Referer required"})
			return
		}
		c.Next()
	}
}

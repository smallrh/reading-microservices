package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"net/http"
	"time"
)

func RateLimit(rdb *redis.Client, requests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 跳过健康检查等端点
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/healthz" {
			c.Next()
			return
		}

		clientIP := c.ClientIP()
		if clientIP == "" {
			clientIP = "unknown"
		}

		key := "rate_limit:" + clientIP + ":" + c.Request.URL.Path

		// 使用管道提高性能
		pipe := rdb.Pipeline()
		incr := pipe.Incr(c, key)
		pipe.Expire(c, key, window)

		_, err := pipe.Exec(c)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		count := incr.Val()
		if count > int64(requests) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":       "Too many requests",
				"retry_after": window.Seconds(),
			})
			return
		}

		c.Next()
	}
}

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Next()
	}
}

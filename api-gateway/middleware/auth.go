package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reading-microservices/shared/utils"
	"strings"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Authorization required"})
			c.Abort()
			return
		}
		if len(token) > 7 && strings.ToUpper(token[:7]) == "BEARER " {
			token = token[7:]
		}
		claims, err := utils.ParseToken(token, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Request.Header.Set("X-User-ID", claims.UserID)
		c.Request.Header.Set("X-Username", claims.Username)
		c.Next()
	}
}

func OptionalAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.Next()
			return
		}
		if len(token) > 7 && strings.ToUpper(token[:7]) == "BEARER " {
			token = token[7:]
		}
		claims, err := utils.ParseToken(token, jwtSecret)
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Request.Header.Set("X-User-ID", claims.UserID)
			c.Request.Header.Set("X-Username", claims.Username)
		}
		c.Next()
	}
}

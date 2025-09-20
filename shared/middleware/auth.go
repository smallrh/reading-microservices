package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"reading-microservices/shared/utils"
	"strings"
)

func JWTAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			utils.ErrorWithCode(c, utils.ERROR_UNAUTHORIZED)
			c.Abort()
			return
		}

		// Remove Bearer prefix
		if len(token) > 7 && strings.ToUpper(token[:7]) == "BEARER " {
			token = token[7:]
		}

		claims, err := utils.ParseToken(token, secret)
		if err != nil {
			utils.Error(c, utils.ERROR_UNAUTHORIZED, "Invalid token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Cache-Control, X-File-Name")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}

		c.Next()
	}
}

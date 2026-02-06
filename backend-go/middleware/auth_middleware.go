package middleware

import (
	"net/http"
	"smartfarm-api/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Cek Cookie
		cookie, err := c.Cookie("access_token")
		if err == nil {
			tokenString = cookie
		}

		// 2. Cek Header Authorization (Bearer ...)
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				split := strings.Split(authHeader, " ")
				if len(split) == 2 {
					tokenString = split[1]
				}
			}
		}

		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Set context
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

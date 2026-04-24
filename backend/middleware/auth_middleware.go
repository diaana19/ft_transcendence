package middleware

import (
	"net/http"
	"strings"

	"github.com/Transcendence/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware rejects requests that do not carry a valid Bearer token.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserId)
		c.Next()
	}
}

// OptionalAuthMiddleware parses the token when present but never rejects the request.
// Controllers can check whether "user_id" is set to personalise the response (e.g. liked=true).
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			if claims, err := utils.ValidateJWT(tokenStr); err == nil {
				c.Set("user_id", claims.UserId)
			}
		}
		c.Next()
	}
}

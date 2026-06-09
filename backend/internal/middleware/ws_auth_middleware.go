package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/utils"
)

// WSAuthMiddleware checks the JWT for WebSocket routes. It reads the token from the
// auth_token cookie or the token query param, and blocks the request when it is not valid.
func WSAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie("auth_token")

		if token == "" {
			token = c.Query("token")
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set("user_id", claims.Subject)
		c.Set("username", claims.Username)
		c.Next()
	}
}

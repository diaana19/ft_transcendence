package middleware

import (
	"net/http"

	"github.com/Transcendence/utils"
	"github.com/gin-gonic/gin"
)

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

		c.Set("user_id", claims.UserId)
		c.Set("username", claims.Username)
		c.Next()
	}
}

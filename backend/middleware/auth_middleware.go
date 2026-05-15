package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Transcendence/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != authHeader {
			return token
		}
	}

	cookieToken, err := c.Cookie("auth_token")
	if err == nil && cookieToken != "" {
		return cookieToken
	}

	return ""
}

func AuthMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing token",
			})
			return
		}

		ctx := context.Background()
		exists, err := rdb.Exists(ctx, "blacklist:"+token).Result()
		if err == nil && exists > 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token revoked",
			})
			return
		}

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		c.Set("userID", claims.UserId)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.Next()
			return
		}

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserId)
		c.Next()
	}
}

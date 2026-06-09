package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"ft_transcendence/backend/internal/utils"
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

func clearAuthCookie(c *gin.Context) {
	if _, err := c.Cookie("auth_token"); err == nil {
		c.SetCookie("auth_token", "", -1, "/", "", true, true)
	}
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
			clearAuthCookie(c)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "token revoked",
			})
			return
		}

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			clearAuthCookie(c)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token",
			})
			return
		}

		c.Set("user_id", claims.Subject)
		c.Set("username", claims.Username)
		c.Next()
	}
}

func OptionalAuthMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			c.Next()
			return
		}

		ctx := context.Background()
		if exists, err := rdb.Exists(ctx, "blacklist:"+token).Result(); err == nil && exists > 0 {
			c.Next()
			return
		}

		claims, err := utils.ValidateJWT(token)
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.Subject)
		c.Set("username", claims.Username)
		c.Next()
	}
}

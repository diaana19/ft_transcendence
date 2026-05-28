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

		c.Set("user_id", claims.Subject)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// OptionalAuthMiddleware accepts requests with or without a token. When a
// token is present it is validated through the same blacklist + signature
// pipeline as AuthMiddleware — so a logged-out user (whose JWT was just
// blacklisted) does not get personalized treatment on public endpoints like
// GET /api/posts. Invalid or revoked tokens silently fall through to the
// anonymous path rather than returning 401, because the caller did not ask
// for auth.
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

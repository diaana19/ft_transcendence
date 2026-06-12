package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/utils"
)

// extractToken reads the JWT from the Authorization Bearer header or the auth_token cookie.
// It returns an empty string when no token is found.
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

// clearAuthCookie removes the auth_token cookie from the client when it exists.
func clearAuthCookie(c *gin.Context) {
	if _, err := c.Cookie("auth_token"); err == nil {
		c.SetCookie("auth_token", "", -1, "/", "", true, true)
	}
}

// userExists returns true if the user behind the token still has an active account.
// A valid JWT can outlive the account (deleted via GDPR or profile), so the token
// alone is not enough.
func userExists(db *gorm.DB, userID string) bool {
	var count int64
	err := db.Model(&models.User{}).Where("id = ?", userID).Count(&count).Error
	return err == nil && count > 0
}

// AuthMiddleware checks the JWT and blocks the request when it is not valid. It also rejects
// tokens that are in the Redis blacklist or whose user no longer exists, and sets user_id
// and username in the context.
func AuthMiddleware(rdb *redis.Client, db *gorm.DB) gin.HandlerFunc {
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

		if !userExists(db, claims.Subject) {
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

// OptionalAuthMiddleware tries to read the JWT but never blocks the request. When the token
// is valid and its user still exists it sets user_id and username in the context, otherwise
// it just continues as an anonymous request.
func OptionalAuthMiddleware(rdb *redis.Client, db *gorm.DB) gin.HandlerFunc {
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

		if !userExists(db, claims.Subject) {
			c.Next()
			return
		}

		c.Set("user_id", claims.Subject)
		c.Set("username", claims.Username)
		c.Next()
	}
}

package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimitMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := fmt.Sprintf("rate_limit:%s", ip)

		count, err := rdb.Incr(c.Request.Context(), key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limiter error"})
			c.Abort()
			return
		}

		if count == 1 {
			rdb.Expire(c.Request.Context(), key, time.Minute)
		}

		if count > 20 {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many request please try again in a minute."})
			c.Abort()
			return
		}
		c.Next()
	}
}

package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const defaultRateLimitMax = 20

func rateLimitMax() int64 {
	if v := os.Getenv("RATE_LIMIT_MAX"); v != "" {
		if n, err := strconv.ParseInt(v, 10, 64); err == nil {
			return n
		}
	}
	return defaultRateLimitMax
}

func RateLimitMiddleware(rdb *redis.Client) gin.HandlerFunc {
	maxRequests := rateLimitMax()
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

		if count > maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "too many request please try again in a minute."})
			c.Abort()
			return
		}
		c.Next()
	}
}

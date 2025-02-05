package middleware

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var ctx = context.Background()

// RateLimiter middleware to limit requests
func RateLimiter(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		token := c.GetHeader("API_KEY")

		// Determine the key to use for rate limiting
		key := ip
		if token != "" {
			key = token
		}

		// Get the rate limit and block duration from environment variables
		rateLimit, _ := strconv.Atoi(os.Getenv("RATE_LIMIT"))
		blockDuration, _ := strconv.Atoi(os.Getenv("BLOCK_DURATION"))

		// Increment the request count for the key
		count, err := rdb.Incr(ctx, key).Result()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}

		// Set the expiration for the key if it's the first request
		if count == 1 {
			rdb.Expire(ctx, key, time.Second)
		}

		// Check if the request count exceeds the rate limit
		if count > int64(rateLimit) {
			// Set the expiration for the key to block duration
			rdb.Expire(ctx, key, time.Duration(blockDuration)*time.Second)
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "you have reached the maximum number of requests or actions allowed within a certain time frame"})
			c.Abort()
			return
		}

		// Continue to the next handler
		c.Next()
	}
}

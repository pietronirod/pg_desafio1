package limiter

import (
	"net/http"
	"rate-limiter/config"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RateLimiterMiddleware(rateLimiter *RateLimiterService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		token := c.GetHeader("API_KEY")

		result := rateLimiter.AllowRequest(ip, token)
		if !result.Allowed {
			config.Logger.Warn("Requisição bloqueada pelo Rate Limiter", zap.String("ip", ip), zap.String("token", token))

			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message":  "You have reached the maximum number of requests or actions allowed within a certain time frame",
				"retry_in": result.BlockTime.Seconds(),
			})
			return
		}

		c.Next()
	}
}

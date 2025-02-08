package main

import (
	"fmt"
	"rate-limiter/config"
	"rate-limiter/internal/limiter"
	"rate-limiter/internal/storage"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	config.LoadConfig()
	config.InitLogger()

	var rateLimiterStorage storage.RateLimiterStorage
	if config.Cfg.RedisAddr != "" {
		rateLimiterStorage = storage.NewRedisStorage(config.Cfg.RedisAddr, config.Cfg.RedisPassword, config.Cfg.RedisDB)
	} else {
		rateLimiterStorage = storage.NewMemoryStorage()
	}

	rateLimiterService := limiter.NewRateLimiterService(rateLimiterStorage, limiter.RateLimiterConfig{
		RateLimitPerIP:        config.Cfg.RateLimitPerIP,
		RateLimitPerToken:     config.Cfg.RateLimitPerToken,
		BlockTimePerIP:        config.Cfg.BlockTimePerIP,
		BlockTimePerToken:     config.Cfg.BlockTimePerToken,
		DefaultBlockTimeIP:    config.Cfg.DefaultBlockTimeIP,
		DefaultBlockTimeToken: config.Cfg.DefaultBlockTimeToken,
	}, config.Logger)

	r := gin.Default()
	r.Use(limiter.RateLimiterMiddleware(rateLimiterService))

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Requisição permitida"})
	})

	port := config.Cfg.ServerPort
	fmt.Println("Servidor rodando na porta", port)
	if err := r.Run(":" + port); err != nil {
		config.Logger.Fatal("Erro ao iniciar servidor", zap.Error(err))
	}

}

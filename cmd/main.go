package main

import (
	"net/http"
	"os"

	"rate-limiter/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load the .env file
	err := godotenv.Load("cmd/.env")
	if err != nil {
		zap.L().Fatal("Error loading .env file", zap.Error(err))
	}

	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Initialize Gin router
	router := gin.Default()
	router.Use(middleware.RateLimiter(rdb))

	// Define routes
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})

	// Start the server
	zap.L().Info("Starting server on :8080")
	router.Run(":8080")
}

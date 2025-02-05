package main

import (
	"log"
	"net/http"

	"rate-limiter/internal/middleware"
	"rate-limiter/internal/redis"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load("cmd/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Redis client
	redisClient := redis.NewClient()

	// Create a new Gin router
	router := gin.Default()

	// Apply the rate limiter middleware
	router.Use(middleware.RateLimiter(redisClient))

	// Define a simple handler for testing
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	// Start the server on port 8080
	router.Run(":8080")
}

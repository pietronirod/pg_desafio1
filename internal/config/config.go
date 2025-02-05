package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RequestLimit  int
	BlockDuration int
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	requestLimit, err := strconv.Atoi(getEnv("REQUEST_LIMIT", "100"))
	if err != nil {
		return nil, err
	}

	blockDuration, err := strconv.Atoi(getEnv("BLOCK_DURATION", "60"))
	if err != nil {
		return nil, err
	}

	return &Config{
		RequestLimit:  requestLimit,
		BlockDuration: blockDuration,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

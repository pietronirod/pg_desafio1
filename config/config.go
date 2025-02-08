package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	RateLimitPerIP        int
	RateLimitPerToken     int
	DefaultBlockTimeIP    int
	DefaultBlockTimeToken int
	BlockTimePerIP        map[string]int
	BlockTimePerToken     map[string]int
	RedisAddr             string
	RedisPassword         string
	RedisDB               int
	ServerPort            string
	LogLevel              string
}

var Cfg Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("[WARN] Arquivo .env não encontrado, usando variáveis de ambiente")
	}

	Cfg = Config{
		RateLimitPerIP:        getEnvAsInt("RATE_LIMIT_PER_IP", 5),
		RateLimitPerToken:     getEnvAsInt("RATE_LIMIT_PER_TOKEN", 100),
		DefaultBlockTimeIP:    getEnvAsInt("DEFAULT_BLOCK_TIME_IP", 300),
		DefaultBlockTimeToken: getEnvAsInt("DEFAULT_BLOCK_TIME_TOKEN", 300),
		BlockTimePerIP:        parseBlockTimeList(getEnv("BLOCK_TIME_PER_IP", "")),
		BlockTimePerToken:     parseBlockTimeList(getEnv("BLOCK_TIME_PER_TOKEN", "")),
		RedisAddr:             getEnv("REDIS_ADDR", "localhost:6379"),
		RedisDB:               getEnvAsInt("REDIS_DB", 0),
		ServerPort:            getEnv("SERVER_PORT", "8080"),
		LogLevel:              getEnv("LOG_LEVEL", "info"),
	}
}

func parseBlockTimeList(input string) map[string]int {
	result := make(map[string]int)
	if input == "" {
		return result
	}
	pairs := strings.Split(input, ";")
	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) == 2 {
			timeValue, err := strconv.Atoi(parts[1])
			if err == nil {
				result[parts[0]] = timeValue
			}
		}
	}
	return result
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}

package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupRedisContainer(t *testing.T) *redis.Client {
	req := testcontainers.ContainerRequest{
		Image:        "redis:alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}

	redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}

	host, err := redisC.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	port, err := redisC.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port.Port(),
		Password: "",
		DB:       0,
	})

	return rdb
}

func TestRateLimiter(t *testing.T) {
	os.Setenv("RATE_LIMIT", "5")
	os.Setenv("BLOCK_DURATION", "60")

	rdb := setupRedisContainer(t)
	defer rdb.Close()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(RateLimiter(rdb))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, World!"})
	})

	for i := 0; i < 5; i++ {
		req, _ := http.NewRequest("GET", "/", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
		assert.Equal(t, http.StatusOK, resp.Code)
	}

	req, _ := http.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusTooManyRequests, resp.Code)
}

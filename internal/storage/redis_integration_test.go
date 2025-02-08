package storage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupRedisContainer(t *testing.T) (*RedisRateLimiterStorage, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:7.0",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)

	host, err := redisContainer.Host(ctx)
	assert.NoError(t, err)
	port, err := redisContainer.MappedPort(ctx, "6379")
	assert.NoError(t, err)

	redisAddr := fmt.Sprintf("%s:%s", host, port.Port())

	redisStorage := NewRedisStorage(redisAddr, "", 0)

	cleanup := func() {
		redisContainer.Terminate(ctx)
	}

	return redisStorage, cleanup
}

func TestRedis_Ping(t *testing.T) {
	redisStorage, cleanup := setupRedisContainer(t)
	defer cleanup()

	err := redisStorage.Ping()
	assert.NoError(t, err, "O Redis deveria responder ao Ping() corretamente")
}

func TestRedis_IncrementRequest(t *testing.T) {
	redisStorage, cleanup := setupRedisContainer(t)
	defer cleanup()

	key := "test_ip"
	count, err := redisStorage.IncrementRequest(key)
	assert.NoError(t, err)
	assert.Equal(t, 1, count, "A primeira requisição deve ter contagem 1")

	count, err = redisStorage.IncrementRequest(key)
	assert.NoError(t, err)
	assert.Equal(t, 2, count, "A segunda requisição deve ter contagem 2")
}

func TestRedis_BlockKey(t *testing.T) {
	redisStorage, cleanup := setupRedisContainer(t)
	defer cleanup()

	key := "test_blocked_ip"
	redisStorage.BlockKey(key, 2*time.Second)

	blocked, err := redisStorage.IsBlocked(key)
	assert.NoError(t, err)
	assert.True(t, blocked, "A chave deveria estar bloqueada")

	time.Sleep(3 * time.Second)

	blocked, err = redisStorage.IsBlocked(key)
	assert.NoError(t, err)
	assert.False(t, blocked, "A chave deve estar desbloqueada após o tempo de expiração")
}

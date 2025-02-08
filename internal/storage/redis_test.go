package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Configuração para rodar um container do Redis para testes
func setupTestRedis(t *testing.T) (*RedisRateLimiterStorage, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	redisContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)

	// Obtendo a URL do Redis no container
	endpoint, err := redisContainer.Endpoint(ctx, "")
	assert.NoError(t, err)

	// Criando a conexão com o Redis
	redisStorage := NewRedisStorage(endpoint, "", 0)

	// Função para limpar e parar o container após o teste
	cleanup := func() {
		redisContainer.Terminate(ctx)
	}

	return redisStorage, cleanup
}

func TestRedisRateLimiterStorage(t *testing.T) {
	redisStorage, cleanup := setupTestRedis(t)
	defer cleanup()

	// Chaves de teste
	ipKey := "rate_limiter:ip:192.168.1.100"
	tokenKey := "rate_limiter:token:abc123"

	// Teste: IncrementRequest e GetRequestCount
	count, err := redisStorage.IncrementRequest(ipKey)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	count, err = redisStorage.GetRequestCount(ipKey)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Teste: BlockKey e IsBlocked
	err = redisStorage.BlockKey(ipKey, 5*time.Second)
	assert.NoError(t, err)

	blocked, err := redisStorage.IsBlocked(ipKey)
	assert.NoError(t, err)
	assert.True(t, blocked)

	// Teste: SetBlockDuration e GetBlockDuration
	err = redisStorage.SetBlockDuration(ipKey, 10*time.Second)
	assert.NoError(t, err)

	blockDuration, err := redisStorage.GetBlockDuration(ipKey)
	assert.NoError(t, err)
	assert.Equal(t, 10*time.Second, blockDuration)

	// Teste: ResetKey
	err = redisStorage.ResetKey(ipKey)
	assert.NoError(t, err)

	count, err = redisStorage.GetRequestCount(ipKey)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	// Teste: IncrementRequest para um Token
	count, err = redisStorage.IncrementRequest(tokenKey)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Teste: BlockKey para um Token
	err = redisStorage.BlockKey(tokenKey, 3*time.Second)
	assert.NoError(t, err)

	blocked, err = redisStorage.IsBlocked(tokenKey)
	assert.NoError(t, err)
	assert.True(t, blocked)
}

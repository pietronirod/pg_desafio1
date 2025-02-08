package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMemoryRateLimiterStorage(t *testing.T) {
	storage := NewMemoryStorage()

	ipKey := "rate_limiter:ip:192.168.1.100"
	tokenKey := "rate_limiter:token:abc123"

	// Teste: IncrementRequest e GetRequestCount
	count, err := storage.IncrementRequest(ipKey)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	count, err = storage.GetRequestCount(ipKey)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Teste: BlockKey e IsBlocked
	err = storage.BlockKey(ipKey, 5*time.Second)
	assert.NoError(t, err)

	blocked, err := storage.IsBlocked(ipKey)
	assert.NoError(t, err)
	assert.True(t, blocked)

	// Esperar 6 segundos para verificar se o bloqueio expira
	time.Sleep(6 * time.Second)

	blocked, err = storage.IsBlocked(ipKey)
	assert.NoError(t, err)
	assert.False(t, blocked)

	// Teste: SetBlockDuration e GetBlockDuration
	err = storage.SetBlockDuration(ipKey, 10*time.Second)
	assert.NoError(t, err)

	blockDuration, err := storage.GetBlockDuration(ipKey)
	assert.NoError(t, err)
	assert.Equal(t, 10*time.Second, blockDuration)

	// Teste: ResetKey
	err = storage.ResetKey(ipKey)
	assert.NoError(t, err)

	count, err = storage.GetRequestCount(ipKey)
	assert.NoError(t, err)
	assert.Equal(t, 0, count)

	// Teste: IncrementRequest para um Token
	count, err = storage.IncrementRequest(tokenKey)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)

	// Teste: BlockKey para um Token
	err = storage.BlockKey(tokenKey, 3*time.Second)
	assert.NoError(t, err)

	blocked, err = storage.IsBlocked(tokenKey)
	assert.NoError(t, err)
	assert.True(t, blocked)
}

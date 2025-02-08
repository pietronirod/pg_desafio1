package limiter

import (
	"testing"
	"time"

	"rate-limiter/config"
	"rate-limiter/internal/storage"

	"github.com/stretchr/testify/assert"
)

func setupTestRateLimiter() *RateLimiterService {
	config.InitLogger()

	memStorage := storage.NewMemoryStorage()

	testConfig := RateLimiterConfig{
		RateLimitPerIP:        5,
		RateLimitPerToken:     10,
		BlockTimePerIP:        map[string]int{"192.168.1.1": 120},
		BlockTimePerToken:     map[string]int{"token123": 300},
		DefaultBlockTimeIP:    180,
		DefaultBlockTimeToken: 360,
	}

	return NewRateLimiterService(memStorage, testConfig, config.Logger)
}

func TestRateLimiter_IPBlocking(t *testing.T) {
	limiter := setupTestRateLimiter()
	ip := "192.168.1.1"

	// Permitir requisições até o limite
	for i := 0; i < limiter.config.RateLimitPerIP; i++ {
		result := limiter.AllowRequest(ip, "")
		assert.True(t, result.Allowed)
	}

	// Atingiu o limite, deve bloquear
	result := limiter.AllowRequest(ip, "")
	assert.False(t, result.Allowed)
	assert.Equal(t, 120*time.Second, result.BlockTime)

	// Verificar se o IP está bloqueado
	blocked, _ := limiter.storage.IsBlocked(ip)
	assert.True(t, blocked)
}

func TestRateLimiter_TokenBlocking(t *testing.T) {
	limiter := setupTestRateLimiter()
	token := "token123"

	// Permitir requisições até o limite
	for i := 0; i < limiter.config.RateLimitPerToken; i++ {
		result := limiter.AllowRequest("", token)
		assert.True(t, result.Allowed)
	}

	// Atingiu o limite, deve bloquear
	result := limiter.AllowRequest("", token)
	assert.False(t, result.Allowed)
	assert.Equal(t, 300*time.Second, result.BlockTime)

	// Verificar se o Token está bloqueado
	blocked, _ := limiter.storage.IsBlocked(token)
	assert.True(t, blocked)
}

func TestRateLimiter_TokenTakesPriorityOverIP(t *testing.T) {
	limiter := setupTestRateLimiter()
	ip := "192.168.1.1"
	token := "token123"

	// Atingir o limite do IP
	for i := 0; i < limiter.config.RateLimitPerIP; i++ {
		result := limiter.AllowRequest(ip, "")
		assert.True(t, result.Allowed)
	}
	result := limiter.AllowRequest(ip, "")
	assert.False(t, result.Allowed)

	// Mas se um Token válido for enviado, a requisição deve passar até o limite do Token
	for i := 0; i < limiter.config.RateLimitPerToken; i++ {
		result := limiter.AllowRequest(ip, token)
		assert.True(t, result.Allowed)
	}
	result = limiter.AllowRequest(ip, token)
	assert.False(t, result.Allowed)
}

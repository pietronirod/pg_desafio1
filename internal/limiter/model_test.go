package limiter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiterConfigInitialization(t *testing.T) {
	config := RateLimiterConfig{
		RateLimitPerIP:        10,
		RateLimitPerToken:     100,
		BlockTimePerIP:        map[string]int{"192.168.1.1": 300},
		BlockTimePerToken:     map[string]int{"token123": 600},
		DefaultBlockTimeIP:    180,
		DefaultBlockTimeToken: 360,
	}

	assert.Equal(t, 10, config.RateLimitPerIP)
	assert.Equal(t, 100, config.RateLimitPerToken)
	assert.Equal(t, 300, config.BlockTimePerIP["192.168.1.1"])
	assert.Equal(t, 600, config.BlockTimePerToken["token123"])
	assert.Equal(t, 180, config.DefaultBlockTimeIP)
	assert.Equal(t, 360, config.DefaultBlockTimeToken)
}

func TestRateLimitResultInitialization(t *testing.T) {
	result := RateLimitResult{
		Allowed:   false,
		BlockTime: 30 * time.Second,
	}

	assert.False(t, result.Allowed)
	assert.Equal(t, 30*time.Second, result.BlockTime)
}

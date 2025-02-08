package limiter

import "time"

type RateLimiterConfig struct {
	RateLimitPerIP        int
	RateLimitPerToken     int
	BlockTimePerIP        map[string]int
	BlockTimePerToken     map[string]int
	DefaultBlockTimeIP    int
	DefaultBlockTimeToken int
}

type RateLimitResult struct {
	Allowed   bool
	BlockTime time.Duration
}

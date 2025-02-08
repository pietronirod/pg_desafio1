package limiter

import (
	"rate-limiter/internal/storage"
	"time"

	"go.uber.org/zap"
)

type RateLimiterService struct {
	storage storage.RateLimiterStorage
	config  RateLimiterConfig
	logger  *zap.Logger
}

func NewRateLimiterService(storage storage.RateLimiterStorage, cfg RateLimiterConfig, logger *zap.Logger) *RateLimiterService {
	return &RateLimiterService{
		storage: storage,
		config:  cfg,
		logger:  logger,
	}
}

func (rl *RateLimiterService) AllowRequest(ip, token string) RateLimitResult {
	if token != "" {
		if blocked, _ := rl.storage.IsBlocked(token); blocked {
			blockTime, _ := rl.storage.GetBlockDuration(token)
			rl.logger.Warn("Token bloqueado", zap.String("token", token), zap.Duration("block_time", blockTime))
			return RateLimitResult{Allowed: false, BlockTime: blockTime}
		}

		limit := rl.config.RateLimitPerToken
		requests, _ := rl.storage.IncrementRequest(token)

		if requests > limit {
			blockTime := rl.getBlockDurationForToken(token)
			rl.storage.BlockKey(token, blockTime)
			rl.logger.Warn("Token atingiu o limite", zap.String("token", token), zap.Int("requests", requests))
			return RateLimitResult{Allowed: false, BlockTime: blockTime}
		}

		return RateLimitResult{Allowed: true}
	}

	if blocked, _ := rl.storage.IsBlocked(ip); blocked {
		blockTime, _ := rl.storage.GetBlockDuration(ip)
		rl.logger.Warn("IP bloqueado", zap.String("ip", ip), zap.Duration("block_time", blockTime))
		return RateLimitResult{Allowed: false, BlockTime: blockTime}
	}

	limit := rl.config.RateLimitPerIP
	requests, _ := rl.storage.IncrementRequest(ip)

	if requests > limit {
		blockTime := rl.getBlockDurationForIP(ip)
		rl.storage.BlockKey(ip, blockTime)
		rl.logger.Warn("IP atingiu o limite", zap.String("ip", ip), zap.Int("requests", requests))
		return RateLimitResult{Allowed: false, BlockTime: blockTime}
	}

	return RateLimitResult{Allowed: true}
}

func (rl *RateLimiterService) getBlockDurationForIP(ip string) time.Duration {
	if duration, exists := rl.config.BlockTimePerIP[ip]; exists {
		return time.Duration(duration) * time.Second
	}
	return time.Duration(rl.config.DefaultBlockTimeIP) * time.Second
}

func (rl *RateLimiterService) getBlockDurationForToken(token string) time.Duration {
	if duration, exists := rl.config.BlockTimePerToken[token]; exists {
		return time.Duration(duration) * time.Second
	}
	return time.Duration(rl.config.DefaultBlockTimeToken) * time.Second
}

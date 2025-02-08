package storage

import "time"

type RateLimiterStorage interface {
	IncrementRequest(key string) (int, error)
	GetRequestCount(key string) (int, error)
	BlockKey(key string, duration time.Duration) error
	IsBlocked(key string) (bool, error)
	ResetKey(key string) error
	GetBlockDuration(key string) (time.Duration, error)
	SetBlockDuration(key string, duration time.Duration) error
}

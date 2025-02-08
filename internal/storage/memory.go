package storage

import (
	"sync"
	"time"
)

type MemoryRateLimiterStorage struct {
	mu             sync.Mutex
	requests       map[string]int
	blocked        map[string]time.Time
	blockDurations map[string]time.Duration
}

func NewMemoryStorage() *MemoryRateLimiterStorage {
	return &MemoryRateLimiterStorage{
		requests:       make(map[string]int),
		blocked:        make(map[string]time.Time),
		blockDurations: make(map[string]time.Duration),
	}
}

func (m *MemoryRateLimiterStorage) IncrementRequest(key string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.requests[key]++
	return m.requests[key], nil
}

func (m *MemoryRateLimiterStorage) GetRequestCount(key string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	count, exists := m.requests[key]
	if !exists {
		return 0, nil
	}
	return count, nil
}

func (m *MemoryRateLimiterStorage) BlockKey(key string, duration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.blocked[key] = time.Now().Add(duration)
	m.blockDurations[key] = duration
	return nil
}

func (m *MemoryRateLimiterStorage) IsBlocked(key string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	expiryTime, exists := m.blocked[key]
	if !exists {
		return false, nil
	}
	if time.Now().After(expiryTime) {
		delete(m.blocked, key)
		delete(m.blockDurations, key)
		return false, nil
	}
	return true, nil
}

func (m *MemoryRateLimiterStorage) ResetKey(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.requests, key)
	return nil
}

func (m *MemoryRateLimiterStorage) GetBlockDuration(key string) (time.Duration, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	duration, exists := m.blockDurations[key]
	if !exists {
		return 0, nil
	}
	return duration, nil
}

func (m *MemoryRateLimiterStorage) SetBlockDuration(key string, duration time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.blockDurations[key] = duration
	return nil
}

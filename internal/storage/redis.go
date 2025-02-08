package storage

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRateLimiterStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStorage(redisAddr, redisPassword string, redisDB int) *RedisRateLimiterStorage {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       redisDB,
	})
	return &RedisRateLimiterStorage{
		client: rdb,
		ctx:    context.Background(),
	}
}

func (r *RedisRateLimiterStorage) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := r.client.Ping(ctx).Result()
	return err
}

func (r *RedisRateLimiterStorage) IncrementRequest(key string) (int, error) {
	count, err := r.client.Incr(r.ctx, key).Result()
	if err != nil {
		return 0, err
	}
	r.client.Expire(r.ctx, key, time.Minute)
	return int(count), nil
}

func (r *RedisRateLimiterStorage) GetRequestCount(key string) (int, error) {
	count, err := r.client.Get(r.ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	parsedCount, _ := strconv.Atoi(count)
	return parsedCount, nil
}

func (r *RedisRateLimiterStorage) BlockKey(key string, duration time.Duration) error {
	return r.client.Set(r.ctx, "rate_limiter:block:"+key, 1, duration).Err()
}

func (r *RedisRateLimiterStorage) IsBlocked(key string) (bool, error) {
	exists, err := r.client.Exists(r.ctx, "rate_limiter:block:"+key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}

func (r *RedisRateLimiterStorage) ResetKey(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisRateLimiterStorage) GetBlockDuration(key string) (time.Duration, error) {
	duration, err := r.client.Get(r.ctx, "rate_limiter:block:"+key).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	blockTime, _ := strconv.Atoi(duration)
	return time.Duration(blockTime) * time.Second, nil
}

func (r *RedisRateLimiterStorage) SetBlockDuration(key string, duration time.Duration) error {
	return r.client.Set(r.ctx, "rate_limiter:block:"+key, int(duration.Seconds()), duration).Err()
}

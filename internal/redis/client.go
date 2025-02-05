package redis

import (
	"context"
	"os"

	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisClient(addr string, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisClient{
		client: rdb,
		ctx:    context.Background(),
	}
}

// NewClient initializes a new Redis client
func NewClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return rdb
}

// GetClient returns the Redis client
func GetClient() *redis.Client {
	return NewClient()
}

func (r *RedisClient) Increment(key string) (int64, error) {
	return r.client.Incr(r.ctx, key).Result()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

func (r *RedisClient) Set(key string, value string, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

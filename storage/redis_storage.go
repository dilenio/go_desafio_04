package storage

import (
	"time"

	"github.com/go-redis/redis"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisStorage(addr, password string, db int) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisStorage{client: client}, nil
}

func (rs *RedisStorage) Set(key, value string, timeout time.Duration) error {
	return rs.client.Set(key, value, timeout).Err()
}

func (rs *RedisStorage) Get(key string) (string, error) {
	return rs.client.Get(key).Result()
}

func (rs *RedisStorage) Incr(key string) error {
	return rs.client.Incr(key).Err()
}

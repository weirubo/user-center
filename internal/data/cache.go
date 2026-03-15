package data

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

func NewCache(cfg interface{}) (*Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect redis: %v", err)
	}
	return &Cache{client: client}, nil
}

func (c *Cache) Close() error {
	return c.client.Close()
}

func (c *Cache) Set(key string, value interface{}) error {
	return c.client.Set(context.Background(), key, value, 0).Err()
}

func (c *Cache) Get(key string) (string, error) {
	return c.client.Get(context.Background(), key).Result()
}

func (c *Cache) SetWithExpire(key string, value interface{}, expire time.Duration) error {
	return c.client.Set(context.Background(), key, value, expire).Err()
}

func (c *Cache) Delete(key string) error {
	return c.client.Del(context.Background(), key).Err()
}

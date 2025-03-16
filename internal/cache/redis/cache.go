package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Stuhub-io/logger"
	redis "github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

func Must(url string, logger logger.Logger) *RedisCache {
	opts, err := redis.ParseURL(url)
	if err != nil {
		fmt.Println(">>>> ", err)
		panic("Redis url is not valid!")
	}

	client := redis.NewClient(opts)

	// _, err = client.Ping(context.Background()).Result()
	// if err != nil {
	// 	panic("Make sure you have a Redis server running on the specified host and port")
	// }

	logger.Info("redis connected")

	return &RedisCache{client: client}
}

func (c *RedisCache) Set(key string, value any, duration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value for key %q: %v", key, err)
	}

	_, err = c.client.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("%v", err)
		return nil
	}

	if err := c.client.Set(context.Background(), key, data, duration).Err(); err != nil {
		return fmt.Errorf("failed to set value for key %q: %v", key, err)
	}

	return nil
}

func (c *RedisCache) Get(key string) (string, error) {
	data, err := c.client.Get(context.Background(), key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("cache miss for key %q", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to get value for key %q: %v", key, err)
	}

	return data, nil
}

func (c *RedisCache) Delete(key string) error {
	if err := c.client.Del(context.Background(), key).Err(); err != nil {
		return fmt.Errorf("failed to delete value for key %q: %v", key, err)
	}

	return nil
}

func (c *RedisCache) GetClient() *redis.Client {
	return c.client
}

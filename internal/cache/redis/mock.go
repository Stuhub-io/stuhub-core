package redis

import (
	"fmt"
	"time"

	"github.com/Stuhub-io/core/ports"
)

type MockRedisCache struct {
}

func Mock() ports.Cache {
	return &MockRedisCache{}
}

func (c *MockRedisCache) Set(key string, value any, duration time.Duration) error {
	return nil
}

func (c *MockRedisCache) Get(key string) (string, error) {
	return "", fmt.Errorf("cache miss for key %q", key)
}

func (c *MockRedisCache) Delete(key string) error {
	return nil
}

package ports

import (
	"time"

	"github.com/Stuhub-io/core/domain"
)

type Cache interface {
	Set(key string, value any, duration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

type CacheStore interface {
	SetUser(user *domain.User, duration time.Duration) error
	GetUser(userPkID int64) *domain.User
}

package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Stuhub-io/core/domain"
)

func (u *CacheStore) SetUser(user *domain.User, duration time.Duration) error {
	err := u.cache.Set(domain.UserKey(user.PkID), user, duration)
	if err != nil {
		fmt.Printf("caching user error: %v", err)
	}

	return err
}

func (u *CacheStore) GetUser(userPkID int64) *domain.User {
	var user domain.User

	data, err := u.cache.Get(domain.UserKey(userPkID))
	if err != nil {
		return nil
	}

	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return nil
	}

	return &user
}

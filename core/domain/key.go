package domain

import "fmt"

var (
	UserKey = func(userPkID int64) string { return fmt.Sprintf("user:%d", userPkID) }
)

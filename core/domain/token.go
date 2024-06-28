package domain

import "time"

const (
	AccessTokenDuration  = 24 * time.Hour
	RefreshTokenDuration = 24 * 7 * time.Hour
)

type TokenPayload struct {
	UserPkID  string
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

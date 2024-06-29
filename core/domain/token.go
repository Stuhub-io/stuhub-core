package domain

import "time"

const (
	AccessTokenDuration  = 24 * time.Hour
	RefreshTokenDuration = 24 * 7 * time.Hour
)

const (
	EmailVerificationTokenDuration = 10 * time.Minute
)

type TokenPayload struct {
	UserPkID  int64     `json:"user_pkid"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

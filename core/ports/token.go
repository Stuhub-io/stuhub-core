package ports

import (
	"time"
)

type TokenPayload struct {
	ID        string
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type TokenMaker interface {
	CreateToken(email string, duration time.Duration) (string, error)
	VerifyToken(token string) (*TokenPayload, error)
}

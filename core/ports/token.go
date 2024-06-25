package ports

import (
	"time"

	"github.com/Stuhub-io/core/domain"
)

type TokenMaker interface {
	CreateToken(email string, duration time.Duration) (string, error)
	VerifyToken(token string) (*domain.TokenPayload, error)
}

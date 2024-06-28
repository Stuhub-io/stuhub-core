package ports

import (
	"time"

	"github.com/Stuhub-io/core/domain"
)

type TokenMaker interface {
	CreateToken(id string, email string, duration time.Duration) (string, error)
	DecodeToken(token string) (*domain.TokenPayload, error)
}

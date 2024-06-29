package ports

import (
	"time"

	"github.com/Stuhub-io/core/domain"
)

type TokenMaker interface {
	CreateToken(pkid int64, email string, duration time.Duration) (string, error)
	DecodeToken(token string) (*domain.TokenPayload, error)
}

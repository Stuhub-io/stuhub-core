package ports

import (
	"context"

	"github.com/Stuhub-io/core/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.User, *domain.Error)
	GetByEmail(ctx context.Context, email string) (*domain.User, *domain.Error)
}

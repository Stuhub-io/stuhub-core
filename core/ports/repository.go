package ports

import (
	"context"

	"github.com/Stuhub-io/core/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, *domain.Error)
	GetByEmail(ctx context.Context, email string) (*domain.User, *domain.Error)
	CreateNewUser(ctx context.Context, email string) (*domain.User, *domain.Error)
}

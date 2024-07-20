package ports

import (
	"context"

	"github.com/Stuhub-io/core/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, *domain.Error)
	GetUserByPkID(ctx context.Context, pkID int64) (*domain.User, *domain.Error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, *domain.Error)
	GetOrCreateUserByEmail(ctx context.Context, email, salt string) (*domain.User, *domain.Error)
	CreateUserWithGoogleInfo(ctx context.Context, email, salt, firstName, lastName, avatar string) (*domain.User, *domain.Error)
	SetUserPassword(ctx context.Context, PkID int64, hashedPassword string) *domain.Error
	CheckPassword(ctx context.Context, email, rawPassword string, hasher Hasher) (bool, *domain.Error)
	UpdateUserInfo(ctx context.Context, PkID int64, firstName, lastName, avatar string) (*domain.User, *domain.Error)
}

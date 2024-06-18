package postgres

import (
	"context"

	"github.com/Stuhub-io/core/domain"
	"github.com/uptrace/bun"
)

type UserRepository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	// handle DB
	return nil, nil
}

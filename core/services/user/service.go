package user

import (
	"context"
	"time"

	"github.com/Stuhub-io/core/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*domain.User, error)
}

type Service struct {
	userRepository UserRepository
}

type NewServiceParams struct {
	UserRepository
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		userRepository: params.UserRepository,
	}
}

func (s *Service) Login(loginDto LoginDto) (*LoginResponse, error) {
	return &LoginResponse{
		User: domain.User{
			ID:        1,
			Username:  "Khoa 2 updated",
			CreatedAt: time.DateTime,
			UpdatedAt: time.DateTime,
		},
		Token: "token",
	}, nil
}

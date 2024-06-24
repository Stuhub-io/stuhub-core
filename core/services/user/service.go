package user

import (
	"context"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
)

type Service struct {
	userRepository ports.UserRepository
}

type NewServiceParams struct {
	ports.UserRepository
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		userRepository: params.UserRepository,
	}
}

func (s *Service) GetUserById(id int64) (*GetUserByIdResponse, *domain.Error) {
	user, err := s.userRepository.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return &GetUserByIdResponse{
		User: user,
	}, nil
}

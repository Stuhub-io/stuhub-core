package user

import (
	"context"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
)

type Service struct {
	userRepository ports.UserRepository
	cfg            config.Config
}

type NewServiceParams struct {
	ports.UserRepository
	config.Config
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		userRepository: params.UserRepository,
		cfg:            params.Config,
	}
}

func (s *Service) GetUserById(id string) (*GetUserByIdResponse, *domain.Error) {
	user, err := s.userRepository.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return &GetUserByIdResponse{
		User: user,
	}, nil
}

func (s *Service) GetUserByEmail(email string) (*GetUserByEmailResponse, *domain.Error) {
	user, err := s.userRepository.GetUserByEmail(context.Background(), email)
	if err != nil && err.Error != domain.NotFoundErr {
		return nil, err
	}

	return &GetUserByEmailResponse{
		User: user,
	}, nil
}

func (s *Service) UpdateUserInfo(pkID int64, firstName, lastName, avatar string) (*UpdateUserInfoResponse, *domain.Error) {
	user, err := s.userRepository.UpdateUserInfo(context.Background(), pkID, firstName, lastName, avatar)

	if err != nil {
		return nil, err
	}

	return &UpdateUserInfoResponse{
		User: user,
	}, nil
}

func (s *Service) SearchUsers(input domain.UserSearchQuery, currentUser *domain.User) ([]domain.User, *domain.Error) {
	users, err := s.userRepository.Search(context.Background(), input, currentUser)
	if err != nil {
		return nil, err
	}

	return users, nil
}

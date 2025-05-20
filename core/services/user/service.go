package user

import (
	"context"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
)

type Service struct {
	repo *ports.Repository
	cfg  config.Config
}

type NewServiceParams struct {
	*ports.Repository
	config.Config
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		repo: params.Repository,
		cfg:  params.Config,
	}
}

func (s *Service) GetUserById(id string) (*GetUserByIdResponse, *domain.Error) {
	user, err := s.repo.User.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	return &GetUserByIdResponse{
		User: user,
	}, nil
}

func (s *Service) GetUserByEmail(email string) (*GetUserByEmailResponse, *domain.Error) {
	user, err := s.repo.User.GetUserByEmail(context.Background(), email)
	if err != nil && err.Error != domain.NotFoundErr {
		return nil, err
	}

	return &GetUserByEmailResponse{
		User: user,
	}, nil
}

func (s *Service) UpdateUserInfo(pkID int64, firstName, lastName, avatar string) (*UpdateUserInfoResponse, *domain.Error) {
	user, err := s.repo.User.UpdateUserInfo(context.Background(), pkID, firstName, lastName, avatar)

	if err != nil {
		return nil, err
	}

	return &UpdateUserInfoResponse{
		User: user,
	}, nil
}

func (s *Service) SearchUsers(input domain.UserSearchQuery, currentUser *domain.User) ([]domain.User, *domain.Error) {
	users, err := s.repo.User.Search(context.Background(), input, currentUser)
	if err != nil {
		return nil, err
	}

	return users, nil
}

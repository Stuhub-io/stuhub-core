package space

import (
	"context"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
)

type Service struct {
	cfg             config.Config
	spaceRepository ports.SpaceRepository
}

type NewServiceParams struct {
	config.Config
	ports.SpaceRepository
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		cfg:             params.Config,
		spaceRepository: params.SpaceRepository,
	}
}

func (s *Service) CreateOrgSpace(dto CreateSpaceDto) (*domain.Space, *domain.Error) {
	org, err := s.spaceRepository.CreateSpace(context.Background(), dto.OrgPkID, dto.OwnerPkID, false, dto.Name, dto.Description)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (s *Service) GetJoinedSpaceByOrgPkID(orgPkID int64, currentUserPkID int64) ([]domain.Space, *domain.Error) {
	spaces, err := s.spaceRepository.GetSpacesByOrgPkID(context.Background(), orgPkID)
	if err != nil {
		return nil, err
	}
	return spaces, nil
}

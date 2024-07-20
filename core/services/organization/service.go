package organization

import (
	"context"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
)

type Service struct {
	cfg           config.Config
	orgRepository ports.OrganizationRepository
}

type NewServiceParams struct {
	config.Config
	ports.OrganizationRepository
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		cfg:           params.Config,
		orgRepository: params.OrganizationRepository,
	}
}

func (s *Service) CreateOrganization(dto CreateOrganizationParams) (*domain.Organization, *domain.Error) {
	org, err := s.orgRepository.CreateOrg(context.Background(), dto.OwnerPkID, dto.Name, dto.Description, dto.Avatar)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (s *Service) GetOrganizationDetailBySlug(slug string) (*domain.Organization, *domain.Error) {
	org, err := s.orgRepository.GetOrgBySlug(context.Background(), slug)
	if err != nil {
		return nil, err
	}
	return org, nil
}

func (s *Service) GetJoinedOrgs(userPkID int64) ([]domain.Organization, *domain.Error) {
	orgs, err := s.orgRepository.GetOrgsByUserPkID(context.Background(), userPkID)
	if err != nil {
		return nil, err
	}
	// FIXME: Implement sort by recent visit later
	return orgs, nil
}

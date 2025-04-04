package activity

import (
	"context"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/logger"
)

type Service struct {
	cfg                config.Config
	logger             logger.Logger
	pageRepository     ports.PageRepository
	activityRepository ports.ActivityRepository
	orgRepository      ports.OrganizationRepository
}

type NewServiceParams struct {
	config.Config
	logger.Logger
	ports.PageRepository
	ports.ActivityRepository
	ports.OrganizationRepository
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		cfg:                params.Config,
		logger:             params.Logger,
		pageRepository:     params.PageRepository,
		activityRepository: params.ActivityRepository,
		orgRepository:      params.OrganizationRepository,
	}
}

func (s Service) TrackUserVisitPage(curUser *domain.User, pagePkID int64) *domain.Error {
	if curUser == nil {
		return domain.ErrUnauthorized
	}

	p, e := s.pageRepository.GetByID(context.Background(), "", &pagePkID, domain.PageDetailOptions{}, nil)

	if e != nil {
		return e
	}

	label := "User Visited Page"
	_, er := s.activityRepository.Create(context.Background(), domain.ActivityInput{
		ActionCode: domain.ActionUserVisitPage,
		ActorPkID:  curUser.PkID,
		PagePkID:   &p.PkID,
		OrgPkID:    &p.OrganizationPkID,
		Label:      &label,
	})
	if er != nil {
		return er
	}
	return nil
}

func (s Service) TrackUserVisitOrganization(curUser *domain.User, orgPkID int64) *domain.Error {
	if curUser == nil {
		return domain.ErrUnauthorized
	}
	org, e := s.orgRepository.GetOrgByPkID(context.Background(), orgPkID)

	if e != nil {
		return e
	}

	label := "User Visited Page"

	_, er := s.activityRepository.Create(context.Background(), domain.ActivityInput{
		ActionCode: domain.ActionUserVisitPage,
		ActorPkID:  curUser.PkID,
		PagePkID:   nil,
		OrgPkID:    &org.PkId,
		Label:      &label,
	})
	if er != nil {
		return er
	}

	return nil
}

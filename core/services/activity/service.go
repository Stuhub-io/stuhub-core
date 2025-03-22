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
	activityRepository ports.ActivityRepository
	pageRepository     ports.PageRepository
	mailer             ports.Mailer
}

type NewServiceParams struct {
	config.Config
	logger.Logger
	ports.ActivityRepository
	ports.PageRepository
	ports.Mailer
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		cfg:                params.Config,
		logger:             params.Logger,
		activityRepository: params.ActivityRepository,
		pageRepository:     params.PageRepository,
		mailer:             params.Mailer,
	}
}

func (s *Service) CreateActivity(input domain.ActivityInput) (d *domain.Activity, e *domain.Error) {
	d, e = s.activityRepository.Create(context.Background(), input)
	return d, e
}

func (s *Service) ListActivity(query domain.ActivityListQuery) (d []domain.Activity, e *domain.Error) {
	d, e = s.activityRepository.List(context.Background(), query)
	return d, e
}

func (s *Service) ListLatestActivity(query domain.ActivityListQuery) (d []domain.Activity, e *domain.Error) {
	return d, e
}

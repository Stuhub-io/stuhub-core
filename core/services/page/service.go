package page

import (
	"context"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
)

type Service struct {
	cfg            config.Config
	pageRepository ports.PageRepository
	docRepository  ports.DocumentRepository
}

type NewServiceParams struct {
	config.Config
	ports.PageRepository
	ports.DocumentRepository
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		cfg:            params.Config,
		pageRepository: params.PageRepository,
		docRepository:  params.DocumentRepository,
	}
}

func (s *Service) CreateNewPage(dto CreatePageDto) (*domain.Page, *domain.Error) {
	page, err := s.pageRepository.CreatePage(context.Background(), dto.SpacePkID, dto.Name, dto.ViewType, dto.ParentPagePkID)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (s *Service) GetPagesBySpacePkID(spacePkID int64) ([]domain.Page, *domain.Error) {
	pages, err := s.pageRepository.GetPagesBySpacePkID(context.Background(), spacePkID)
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (s *Service) DeletePageByPkID(pagePkID int64, userPkID int64) (*domain.Page, *domain.Error) {
	result, err := s.pageRepository.DeletePageByPkID(context.Background(), pagePkID, userPkID)
	if err != nil {
		return nil, err
	}
	return result, nil
}
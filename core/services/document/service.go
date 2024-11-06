package document

import (
	"context"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
)

type Service struct {
	cfg           config.Config
	docRepository ports.DocumentRepository
}

type NewServiceParams struct {
	config.Config
	ports.DocumentRepository
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		cfg:           params.Config,
		docRepository: params.DocumentRepository,
	}
}

func (s *Service) CreatePage(pageInput domain.PageInput) (*domain.Page, *domain.Error) {
	page, err := s.docRepository.CreatePage(context.Background(), pageInput)
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}
	return page, nil
}

func (s *Service) GetPagesByOrgPkID(query domain.PageListQuery) (d []domain.Page, e *domain.Error) {
	d, e = s.docRepository.List(context.Background(), query)
	return d, e
}

func (s *Service) UpdatePageByPkID(pagePkID int64, updateInput domain.PageUpdateInput) (d *domain.Page, e *domain.Error) {
	d, e = s.docRepository.Update(context.Background(), pagePkID, updateInput)
	return d, e
}

func (s *Service) GetPageDetailByID(pageID string) (d *domain.Page, e *domain.Error) {
	d, e = s.docRepository.GetByID(context.Background(), pageID)
	return d, e
}

func (s *Service) ArchivedPageByPkID(pagePkID int64) (d *domain.Page, e *domain.Error) {
	d, e = s.docRepository.Archive(context.Background(), pagePkID)
	return d, e
}

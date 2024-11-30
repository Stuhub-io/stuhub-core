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
}

type NewServiceParams struct {
	config.Config
	ports.PageRepository
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		cfg:            params.Config,
		pageRepository: params.PageRepository,
	}
}

// Page Controller

func (s *Service) GetPagesByOrgPkID(query domain.PageListQuery) (d []domain.Page, e *domain.Error) {
	d, e = s.pageRepository.List(context.Background(), query)
	return d, e
}

func (s *Service) UpdatePageByPkID(pagePkID int64, updateInput domain.PageUpdateInput) (d *domain.Page, e *domain.Error) {
	d, e = s.pageRepository.Update(context.Background(), pagePkID, updateInput)
	return d, e
}

func (s *Service) GetPageDetailByID(pageID string) (d *domain.Page, e *domain.Error) {
	d, e = s.pageRepository.GetByID(context.Background(), pageID)
	return d, e
}

func (s *Service) ArchivedPageByPkID(pagePkID int64) (d *domain.Page, e *domain.Error) {
	// Recursive archive all children
	d, e = s.pageRepository.Archive(context.Background(), pagePkID)
	return d, e
}
func (s *Service) MovePageByPkID(pagePkID int64, moveInput domain.PageMoveInput) (d *domain.Page, e *domain.Error) {
	d, e = s.pageRepository.Move(context.Background(), pagePkID, moveInput.ParentPagePkID)
	return d, e
}

// Document Controller

func (s *Service) CreateDocumentPage(pageInput domain.DocumentPageInput) (*domain.Page, *domain.Error) {
	page, err := s.pageRepository.CreateDocumentPage(context.Background(), pageInput)
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}
	return page, nil
}

func (s *Service) UpdateDocumentContentByPkID(pagePkID int64, content domain.DocumentInput) (d *domain.Page, e *domain.Error) {
	d, e = s.pageRepository.UpdateContent(context.Background(), pagePkID, content)
	return d, e
}

// Asset Controller

func (s *Service) CreateAssetPage(assetInput domain.AssetPageInput) (*domain.Page, *domain.Error) {
	page, err := s.pageRepository.CreateAsset(context.Background(), assetInput)
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}
	return page, nil
}

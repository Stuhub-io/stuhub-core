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

func (s *Service) GetPageDetailByID(pageID string, tokenID string) (d *domain.Page, e *domain.Error) {
	var PagePkID *int64
	if pageID == "" {
		token, err := s.pageRepository.GetPublicTokenByID(context.Background(), tokenID)
		if token.ArchivedAt != "" {
			return nil, domain.NewErr("Public page is expired", domain.ResourceInvalidOrExpiredCode)
		}
		if err != nil {
			return nil, domain.ErrDatabaseQuery
		}
		PagePkID = &token.PagePkID
	}
	d, e = s.pageRepository.GetByID(context.Background(), pageID, PagePkID)
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

func (s *Service) CreatePublicPageToken(pageID string) (d *domain.PagePublicToken, e *domain.Error) {
	page, err := s.pageRepository.GetByID(context.Background(), pageID, nil)
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}
	d, e = s.pageRepository.CreatePublicToken(context.Background(), page.PkID)
	return d, e
}

func (s *Service) ArchiveAllPublicPageToken(pageID string) (e *domain.Error) {
	page, err := s.pageRepository.GetByID(context.Background(), pageID, nil)
	if err != nil {
		return domain.ErrDatabaseQuery
	}
	e = s.pageRepository.ArchiveAllPublicToken(context.Background(), page.PkID)
	return e
}

// Document Controller.
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

// Generate Public Token.
func (s *Service) GenerateDocumentPublicToken(pagePkID int64) (d string, e *domain.Error) {
	// d, e = s.pageRepository.GeneratePublicToken(context.Background(), pagePkID)
	return d, e
}

func (s *Service) ValidateDocumentPublicToken(token string) (d bool, e *domain.Error) {
	// d, e = s.pageRepository.ValidatePublicToken(context.Background(), token)
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

package page

import (
	"context"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/logger"
)

type Service struct {
	cfg            config.Config
	pageRepository ports.PageRepository
	docRepository  ports.DocumentRepository
	logger         logger.Logger
}

type NewServiceParams struct {
	config.Config
	ports.PageRepository
	ports.DocumentRepository
	logger.Logger
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		cfg:            params.Config,
		pageRepository: params.PageRepository,
		docRepository:  params.DocumentRepository,
		logger:         params.Logger,
	}
}

func (s *Service) CreateNewPage(dto domain.PageInput) (*domain.Page, *domain.Error) {
	page, err := s.pageRepository.CreatePage(context.Background(), dto.SpacePkID, dto.Name, dto.ViewType, dto.ParentPagePkID, dto.NodeID)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (s *Service) GetPagesBySpacePkID(spacePkID int64) ([]domain.Page, *domain.Error) {
	pages, err := s.pageRepository.GetPagesBySpacePkID(context.Background(), spacePkID, false)
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

func (s *Service) GetPageByID(pageID string) (*domain.Page, *domain.Error) {
	page, err := s.pageRepository.GetPageByID(context.Background(), pageID)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (s *Service) UpdatePageById(pageID string, newPage domain.PageInput) (*domain.Page, *domain.Error) {
	page, err := s.pageRepository.UpdatePageByID(context.Background(), pageID, newPage)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (s *Service) ArchivedPageByID(pageID string) (*domain.Page, *domain.Error) {
	page, err := s.pageRepository.ArchivedPageByID(context.Background(), pageID)
	if err != nil {
		return nil, err
	}
	return page, nil
}

func (s *Service) BulkGetOrCreatePageByNodeID(newPagesInput []domain.PageInput) ([]domain.Page, *domain.Error) {
	nodeIds := make([]string, len(newPagesInput))
	for i, page := range newPagesInput {
		nodeIds[i] = page.NodeID
	}
	existedPages, err := s.pageRepository.GetPagesByNodeID(context.Background(), nodeIds)
	if err != nil {
		return nil, err
	}

	if len(existedPages) == len(newPagesInput) {
		return existedPages, nil
	}

	toCreatePageInputs := make([]domain.PageInput, 0, len(newPagesInput)-len(existedPages))
	for _, newPage := range newPagesInput {

		found := false
		for _, existedPage := range existedPages {
			if newPage.NodeID == existedPage.NodeID {
				found = true
				break
			}
		}
		if !found {
			toCreatePageInputs = append(toCreatePageInputs, newPage)
		}
	}
	createdPages, rerr := s.pageRepository.BulkCreatePages(context.Background(), toCreatePageInputs)
	if rerr != nil {
		return nil, rerr
	}

	pages := append(existedPages, createdPages...)
	return pages, nil
}

func (s *Service) BulkArchivePages(pagePkIDs []int64) *domain.Error {
	err := s.pageRepository.BulkArchivePages(context.Background(), pagePkIDs)
	if err != nil {
		return err
	}
	return nil
}

package document

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

func (s *Service) CreateNewDocument(pagePkID int64, jsonContent string) (*domain.Document, *domain.Error) {
	doc, err := s.docRepository.CreateDocument(context.Background(), pagePkID, jsonContent)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *Service) UpdateDocument(docPkID int64, content string) (*domain.Document, *domain.Error) {
	doc, err := s.docRepository.UpdateDocument(context.Background(), docPkID, content)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *Service) GetDocumentByPagePkID(pagePkID int64) (*domain.Document, *domain.Error) {
	doc, err := s.docRepository.GetDocumentByPagePkID(context.Background(), pagePkID)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

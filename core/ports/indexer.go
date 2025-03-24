package ports

import (
	"context"

	"github.com/Stuhub-io/core/domain"
)

type PageIndexer interface {
	Index(ctx context.Context, page domain.IndexedPage) error
	Search(ctx context.Context, args domain.SearchIndexedPageParams) (*[]domain.QuickSearchPage, error)
	Update(ctx context.Context, page domain.IndexedPage) error
	Delete(ctx context.Context, pageID string) error
}

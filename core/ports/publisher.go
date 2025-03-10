package ports

import (
	"context"

	"github.com/Stuhub-io/core/domain"
)

type PageMessageBrokerPublisher interface {
	Created(ctx context.Context, page *domain.Page) error
	Deleted(ctx context.Context, id string) error
	Updated(ctx context.Context, page *domain.Page) error
}

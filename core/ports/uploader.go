package ports

import (
	"github.com/Stuhub-io/core/domain"
)

type Uploader interface {
	SignRequest(params domain.SignUrlInput) (*domain.SignedUrl, *domain.Error)
}

package upload

import (
	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
)

type UploadService struct {
	cfg      config.Config
	uploader ports.Uploader
}

type NewUploadServiceParams struct {
	Config   config.Config
	Uploader ports.Uploader
}

func NewUploadService(params NewUploadServiceParams) *UploadService {
	return &UploadService{
		cfg:      params.Config,
		uploader: params.Uploader,
	}
}

func (s *UploadService) GenerateSignedUrl(input domain.SignUrlInput) (d *domain.SignedUrl, e *domain.Error) {
	return s.uploader.SignRequest(input)
}

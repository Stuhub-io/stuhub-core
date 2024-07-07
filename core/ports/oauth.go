package ports

import (
	"context"

	"github.com/Stuhub-io/core/domain"
)

type OauthService interface {
	GetGoogleUserInfo(ctx context.Context, token string) (*domain.GoogleUserInfo, error)
}

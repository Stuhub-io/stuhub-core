package oauth

import "github.com/Stuhub-io/logger"

type OauthService struct {
	logger logger.Logger
}

func NewOauthService(logger logger.Logger) *OauthService {
	return &OauthService{
		logger: logger,
	}
}

package oauth

import (
	"context"

	"github.com/Stuhub-io/core/domain"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

func (o *OauthService) GetGoogleUserInfo(ctx context.Context, token string) (*domain.GoogleUserInfo, error) {
	service, err := oauth2.NewService(ctx, option.WithoutAuthentication())
	if err != nil {
		o.logger.Error(err, err.Error())
		return nil, err
	}

	userInfo, err := service.Userinfo.Get().Do(googleapi.QueryParameter("access_token", token))
	if err != nil {
		e, _ := err.(*googleapi.Error)
		o.logger.Error(e, e.Message)
		return nil, e
	}

	return &domain.GoogleUserInfo{
		Email:     userInfo.Email,
		FirstName: userInfo.FamilyName,
		LastName:  userInfo.GivenName,
		Avatar:    userInfo.Picture,
	}, nil
}

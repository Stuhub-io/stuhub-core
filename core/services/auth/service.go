package auth

import (
	"context"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/utils"
)

type Service struct {
	userRepository ports.UserRepository
	mailer         ports.Mailer
	tokenMaker     ports.TokenMaker
	remoteRoute    ports.RemoteRoute
	config         config.Config
}

type NewServiceParams struct {
	ports.UserRepository
	ports.Mailer
	ports.TokenMaker
	ports.RemoteRoute
	config.Config
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		userRepository: params.UserRepository,
		mailer:         params.Mailer,
		tokenMaker:     params.TokenMaker,
		config:         params.Config,
		remoteRoute:    params.RemoteRoute,
	}
}

// Send Magic Link if User not set password
func (s *Service) AuthenByEmailStepOne(dto AuthenByEmailStepOneDto) (*AuthenByEmailStepOneResp, *domain.Error) {
	email := dto.Email
	user, err := s.userRepository.GetOrCreateUserByEmail(context.Background(), email)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	// User can auth with Password
	if user.Password != "" {
		return &AuthenByEmailStepOneResp{
			Email:           user.Email,
			IsRequiredEmail: false,
		}, nil
	}

	token, errToken := s.tokenMaker.CreateToken(user.ID, user.Email, domain.AccessTokenDuration)
	if errToken != nil {
		return nil, domain.ErrInternalServerError
	}

	var url string
	if user.OauthGmail == "" {
		// Send Magic Link For Validate Email and Set Password
		url = s.MakeSetPasswordUrl(token)

	} else {
		// Send Magic Link For Validate Email and Redirect to Oauth
		url = s.MakeValidateOauthUrl(token)
	}
	s.mailer.SendMail(ports.SendSendGridMailPayload{
		FromName:    "Stuhub.IO",
		FromAddress: s.config.SendgridEmailFrom,
		ToName:      utils.GetUserFullName(user.FirstName, user.LastName),
		ToAddress:   user.Email,
		TemplateId:  s.config.SendgridSetPasswordTemplateId,
		Data: map[string]string{
			"link": url,
		},
		Subject: "Authenticate your email",
	})
	return &AuthenByEmailStepOneResp{
		Email:           user.Email,
		IsRequiredEmail: true,
	}, nil
	// Send Magic Link with Oauth redirect
}

func (s *Service) MakeValidateOauthUrl(token string) string {
	baseUrl := s.config.BaseUrl + s.remoteRoute.ValidateEmailOauth
	return baseUrl + "?token=" + token
}

func (s *Service) MakeSetPasswordUrl(token string) string {
	baseUrl := s.config.BaseUrl + s.remoteRoute.SetPassword
	return baseUrl + "?token=" + token
}

// FIXME: return token
func (s *Service) AuthenByEmailPassword(dto AuthenByEmailPassword) *domain.Error {
	return nil
}

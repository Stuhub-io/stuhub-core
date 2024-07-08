package auth

import (
	"context"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/utils/userutils"
)

type Service struct {
	userRepository ports.UserRepository
	oauthService   ports.OauthService
	mailer         ports.Mailer
	tokenMaker     ports.TokenMaker
	remoteRoute    ports.RemoteRoute
	hasher         ports.Hasher
	config         config.Config
}

type NewServiceParams struct {
	ports.UserRepository
	ports.OauthService
	ports.Mailer
	ports.TokenMaker
	ports.RemoteRoute
	ports.Hasher
	config.Config
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		userRepository: params.UserRepository,
		oauthService:   params.OauthService,
		mailer:         params.Mailer,
		tokenMaker:     params.TokenMaker,
		config:         params.Config,
		remoteRoute:    params.RemoteRoute,
		hasher:         params.Hasher,
	}
}

// Send Magic Link if User not set password
func (s *Service) AuthenByEmailStepOne(dto AuthenByEmailStepOneDto) (*AuthenByEmailStepOneResp, *domain.Error) {
	email := dto.Email
	user, err := s.userRepository.GetOrCreateUserByEmail(context.Background(), email, s.hasher.GenerateSalt())
	if err != nil {
		return nil, err
	}

	// User can auth with Password
	if user.HavePassword {
		return &AuthenByEmailStepOneResp{
			Email:           user.Email,
			IsRequiredEmail: false,
		}, nil
	}

	token, errToken := s.tokenMaker.CreateToken(user.PkID, user.Email, domain.EmailVerificationTokenDuration)
	if errToken != nil {
		return nil, domain.ErrInternalServerError
	}

	url := s.MakeValidateEmailAuth(token)
	err = s.mailer.SendMail(ports.SendSendGridMailPayload{
		FromName:    "Stuhub.IO",
		FromAddress: s.config.SendgridEmailFrom,
		ToName:      userutils.GetUserFullName(user.FirstName, user.LastName),
		ToAddress:   user.Email,
		TemplateId:  s.config.SendgridSetPasswordTemplateId,
		Data: map[string]string{
			"url": url,
		},
		Subject: "Authenticate your email",
	})
	if err != nil {
		return nil, err
	}

	return &AuthenByEmailStepOneResp{
		Email:           user.Email,
		IsRequiredEmail: true,
	}, nil
	// Send Magic Link with Oauth redirect
}

func (s *Service) MakeValidateEmailAuth(token string) string {
	baseUrl := s.config.RemoteBaseURL + s.remoteRoute.ValidateEmailOauth

	return baseUrl + "?token=" + token
}

// FIXME: return token
func (s *Service) ValidateEmailAuth(token string) (*ValidateEmailTokenResp, *domain.Error) {
	payload, err := s.tokenMaker.DecodeToken(token)
	if err != nil {
		return nil, domain.ErrTokenExpired
	}

	user, uErr := s.userRepository.GetUserByPkID(context.Background(), payload.UserPkID)
	if uErr != nil {
		return nil, domain.ErrBadRequest
	}

	var providerName string = ""
	if user.OauthGmail != "" {
		providerName = domain.GoogleAuthProvider.Name
	}

	action_token, err := s.tokenMaker.CreateToken(user.PkID, user.Email, domain.NextStepTokenDuration)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	return &ValidateEmailTokenResp{
		Email:        user.Email,
		OAuthPvodier: providerName,
		ActionToken:  action_token,
	}, nil
}

func (s *Service) SetPasswordAndAuthUser(dto AuthenByEmailAfterSetPasswordDto) (*AuthenByEmailStepTwoResp, *domain.Error) {
	user, derr := s.userRepository.GetUserByEmail(context.Background(), dto.Email)
	if derr != nil {
		return nil, domain.ErrUserNotFoundByEmail(dto.Email)
	}

	hashedPassword, err := s.hasher.Hash(dto.RawPassword, user.Salt)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	derr = s.userRepository.SetUserPassword(context.Background(), user.PkID, hashedPassword)
	if derr != nil {
		return nil, domain.ErrInternalServerError
	}

	access, tErr := s.tokenMaker.CreateToken(user.PkID, user.Email, domain.AccessTokenDuration)
	if tErr != nil {
		return nil, domain.ErrInternalServerError
	}

	refresh, tErr := s.tokenMaker.CreateToken(user.PkID, user.Email, domain.RefreshTokenDuration)
	if tErr != nil {
		return nil, domain.ErrInternalServerError
	}

	return &AuthenByEmailStepTwoResp{
		AuthToken: domain.AuthToken{
			Access:  access,
			Refresh: refresh,
		},
	}, nil
}

func (s *Service) AuthenUserByEmailPassword(dto AuthenByEmailPasswordDto) (*domain.AuthToken, *domain.Error) {
	user, derr := s.userRepository.GetUserByEmail(context.Background(), dto.Email)
	if derr != nil {
		return nil, derr
	}

	if !user.HavePassword {
		return nil, domain.ErrBadParamInput
	}

	valid, derr := s.userRepository.CheckPassword(context.Background(), user.Email, dto.RawPassword, s.hasher)
	if derr != nil {
		return nil, domain.ErrInternalServerError
	}

	if !valid {
		return nil, domain.ErrUserPassword
	}

	access, tErr := s.tokenMaker.CreateToken(user.PkID, user.Email, domain.AccessTokenDuration)
	if tErr != nil {
		return nil, domain.ErrInternalServerError
	}

	refresh, tErr := s.tokenMaker.CreateToken(user.PkID, user.Email, domain.RefreshTokenDuration)
	if tErr != nil {
		return nil, domain.ErrInternalServerError
	}

	return &domain.AuthToken{
		Access:  access,
		Refresh: refresh,
	}, nil
}

func (s *Service) AuthenUserByGoogle(dto AuthenByGoogleDto) (*AuthenByGoogleResponse, *domain.Error) {
	userInfo, oerr := s.oauthService.GetGoogleUserInfo(context.Background(), dto.Token)
	if oerr != nil {
		return nil, domain.ErrGetGoogleInfo
	}

	user, err := s.userRepository.GetUserByEmail(context.Background(), userInfo.Email)
	if err != nil && err.Error == domain.NotFoundErr {
		salt := s.hasher.GenerateSalt()
		newUser, err := s.userRepository.CreateUserWithGoogleInfo(context.Background(), userInfo.Email, salt, userInfo.FirstName, userInfo.LastName, userInfo.Avatar)
		if err != nil {
			return nil, err
		}

		user = newUser
	} else if err != nil {
		return nil, err
	}

	access, tErr := s.tokenMaker.CreateToken(user.PkID, user.Email, domain.AccessTokenDuration)
	if tErr != nil {
		return nil, domain.ErrInternalServerError
	}

	refresh, tErr := s.tokenMaker.CreateToken(user.PkID, user.Email, domain.RefreshTokenDuration)
	if tErr != nil {
		return nil, domain.ErrInternalServerError
	}

	return &AuthenByGoogleResponse{
		Profile: user,
		AuthToken: domain.AuthToken{
			Access:  access,
			Refresh: refresh,
		},
	}, nil
}

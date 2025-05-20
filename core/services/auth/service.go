package auth

import (
	"context"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/utils/userutils"
)

type Service struct {
	repo         *ports.Repository
	oauthService ports.OauthService
	mailer       ports.Mailer
	tokenMaker   ports.TokenMaker
	remoteRoute  ports.RemoteRoute
	hasher       ports.Hasher
	config       config.Config
}

type NewServiceParams struct {
	*ports.Repository
	ports.OauthService
	ports.Mailer
	ports.TokenMaker
	ports.RemoteRoute
	ports.Hasher
	config.Config
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		repo:         params.Repository,
		oauthService: params.OauthService,
		mailer:       params.Mailer,
		tokenMaker:   params.TokenMaker,
		config:       params.Config,
		remoteRoute:  params.RemoteRoute,
		hasher:       params.Hasher,
	}
}

// Send Magic Link if User not set password.
func (s *Service) AuthenByEmailStepOne(dto AuthenByEmailStepOneDto) (*AuthenByEmailStepOneResp, *domain.Error) {
	email := dto.Email
	user, err, created := s.repo.User.GetOrCreateUserByEmail(context.Background(), email, s.hasher.GenerateSalt())

	if err != nil {
		return nil, err
	}

	if created {
		go s.repo.Page.SyncPageRoleWithNewUser(context.Background(), *user)
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
		ToName:     userutils.GetUserFullName(user.FirstName, user.LastName),
		ToAddress:  user.Email,
		TemplateId: s.config.SendgridSetPasswordTemplateId,
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

func (s *Service) ValidateEmailAuth(token string) (*ValidateEmailTokenResp, *domain.Error) {
	payload, err := s.tokenMaker.DecodeToken(token)
	if err != nil {
		return nil, domain.ErrTokenExpired
	}

	user, uErr := s.repo.User.GetUserByPkID(context.Background(), payload.UserPkID)
	if uErr != nil {
		return nil, domain.ErrBadRequest
	}

	var providerName string = ""
	if user.OauthGmail != "" {
		providerName = domain.GoogleAuthProvider.Name
	}

	actionToken, err := s.tokenMaker.CreateToken(user.PkID, user.Email, domain.NextStepTokenDuration)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	return &ValidateEmailTokenResp{
		Email:        user.Email,
		OAuthPvodier: providerName,
		ActionToken:  actionToken,
	}, nil
}

func (s *Service) SetPasswordAndAuthUser(dto AuthenByEmailAfterSetPasswordDto) (*AuthenByEmailStepTwoResp, *domain.Error) {
	user, err := s.repo.User.GetUserByEmail(context.Background(), dto.Email)
	if err != nil {
		return nil, domain.ErrUserNotFoundByEmail(dto.Email)
	}

	hashedPassword, herr := s.hasher.Hash(dto.RawPassword, user.Salt)
	if herr != nil {
		return nil, domain.ErrInternalServerError
	}

	err = s.repo.User.SetUserPassword(context.Background(), user.PkID, hashedPassword)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.User.SetUserActivatedAt(context.Background(), user.PkID, time.Now())
	if err != nil {
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

	return &AuthenByEmailStepTwoResp{
		AuthToken: domain.AuthToken{
			Access:  access,
			Refresh: refresh,
		},
	}, nil
}

func (s *Service) ActivateUser(dto ActivateUserDto) (*domain.User, *domain.Error) {
	user, err := s.repo.User.GetUserByPkID(context.Background(), dto.UserPkID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	if user.ActivatedAt != "" {
		return user, nil
	}

	updatedUser, err := s.repo.User.SetUserActivatedAt(context.Background(), dto.UserPkID, time.Now())
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *Service) AuthenUserByEmailPassword(dto AuthenByEmailPasswordDto) (*domain.AuthToken, *domain.User, *domain.Error) {
	user, derr := s.repo.User.GetUserByEmail(context.Background(), dto.Email)
	if derr != nil {
		return nil, nil, domain.ErrUserNotFoundByEmail(dto.Email)
	}

	if !user.HavePassword {
		return nil, nil, domain.ErrBadParamInput
	}

	valid, derr := s.repo.User.CheckPassword(context.Background(), user.Email, dto.RawPassword, s.hasher)
	if derr != nil {
		return nil, nil, domain.ErrInternalServerError
	}

	if !valid {
		return nil, nil, domain.ErrUserPassword
	}

	access, tErr := s.tokenMaker.CreateToken(user.PkID, user.Email, domain.AccessTokenDuration)
	if tErr != nil {
		return nil, nil, domain.ErrInternalServerError
	}

	refresh, tErr := s.tokenMaker.CreateToken(user.PkID, user.Email, domain.RefreshTokenDuration)
	if tErr != nil {
		return nil, nil, domain.ErrInternalServerError
	}

	return &domain.AuthToken{
		Access:  access,
		Refresh: refresh,
	}, user, nil
}

func (s *Service) GetUserByToken(token string) (*domain.User, *domain.Error) {
	payload, err := s.tokenMaker.DecodeToken(token)
	if err != nil {
		return nil, domain.ErrTokenExpired
	}

	cachedUser := s.repo.Store.Cache().GetUser(payload.UserPkID)
	if cachedUser != nil {
		return cachedUser, nil
	}
	user, uErr := s.repo.User.GetUserByPkID(context.Background(), payload.UserPkID)
	if uErr != nil {
		return nil, domain.ErrBadRequest
	}

	return user, nil
}

func (s *Service) AuthenUserByGoogle(dto AuthenByGoogleDto) (*AuthenByGoogleResponse, *domain.Error) {
	userInfo, oerr := s.oauthService.GetGoogleUserInfo(context.Background(), dto.Token)
	if oerr != nil {
		return nil, domain.ErrGetGoogleInfo
	}

	user, err := s.repo.User.GetUserByEmail(context.Background(), userInfo.Email)
	if err != nil && err.Error == domain.NotFoundErr {
		salt := s.hasher.GenerateSalt()
		newUser, err := s.repo.User.CreateUserWithGoogleInfo(context.Background(), userInfo.Email, salt, userInfo.FirstName, userInfo.LastName, userInfo.Avatar)

		if err != nil {
			return nil, err
		}

		go s.repo.Page.SyncPageRoleWithNewUser(context.Background(), *newUser)

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

	if user.ActivatedAt == "" {
		s.repo.User.SetUserActivatedAt(context.Background(), user.PkID, time.Now())
	}

	return &AuthenByGoogleResponse{
		Profile: user,
		AuthToken: domain.AuthToken{
			Access:  access,
			Refresh: refresh,
		},
	}, nil
}

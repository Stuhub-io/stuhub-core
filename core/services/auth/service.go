package auth

import (
	"context"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
)

type Service struct {
	userRepository ports.UserRepository
	mailer         ports.Mailer
	tokenMaker     ports.TokenMaker
}

type NewServiceParams struct {
	ports.UserRepository
	ports.Mailer
	ports.TokenMaker
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		userRepository: params.UserRepository,
		mailer:         params.Mailer,
		tokenMaker:     params.TokenMaker,
	}
}

func (s *Service) RegisterByEmail(registerByEmailDto RegisterByEmailDto) *domain.Error {
	email := registerByEmailDto.Email
	user, _ := s.userRepository.GetByEmail(context.Background(), email)
	if user != nil {
		return domain.ErrExistUserEmail(email)
	}

	//create token & assign to magic link

	//send magic links via mail
	/*
		err := s.mailer.SendMail(ports.SendMailPayload{
			To:          "",
			Address:     "",
			Subject:     "",
			PlainText:   "",
			HTMLContent: "",
		})
		if err != nil {
			return nil, &domain.Error{
				Code:    domain.ErrInternalServerError.Code,
				Error:   domain.ErrInternalServerError.Error,
				Message: err.Error(),
			}
		}
	*/

	return nil
}

// func (s *Service) VerifyMagicLinkToken

// func (s *Service) ActivateAccount

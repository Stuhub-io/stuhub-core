package auth

import (
	"context"
	"fmt"

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

// FIXME: return token
func (s *Service) AuthenByEmailPassword(dto AuthenByEmailPassword) *domain.Error {
	return nil
}

func (s *Service) AuthenWithEmail(dto AuthenByEmailDto) *domain.Error {
	// NOTE Send Magic Link to Email if user passowrd not set
	email := dto.Email
	var user *domain.User
	u, err := s.userRepository.GetByEmail(context.Background(), email)
	if err != nil {
		return err
	}
	if u != nil {
		user = u
	} else {
		u, err := s.userRepository.CreateNewUser(context.Background(), email)
		if err != nil {
			return err
		}
		user = u
	}
	fmt.Print(user)

	// Auth using Password If user already exists

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

// AuthenWithEmail -> send magic link via Email -> verify magic link token
// FIXME: check password set ->
func (s *Service) VerifyMagicLinkToken() *domain.Error {
	return nil
}

//...

// func (s *Service) ActivateAccount

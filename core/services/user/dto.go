package user

import "github.com/Stuhub-io/core/domain"

type GetUserByIdResponse struct {
	User *domain.User `json:"user"`
}

type GetUserByEmailResponse struct {
	User *domain.User `json:"user"`
}

package auth

import "github.com/Stuhub-io/core/domain"

type LoginDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User  domain.User `json:"user"`
	Token string      `json:"token"`
}

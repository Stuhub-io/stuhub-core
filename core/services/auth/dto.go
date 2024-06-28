package auth

import "github.com/Stuhub-io/core/domain"

type AuthenByEmailStepOneDto struct {
	Email string `json:"email"`
}

type AuthenByEmailStepOneResp struct {
	Email           string `json:"email"`
	IsRequiredEmail bool   `json:"is_required_email"`
}

type AuthenByEmailStepTwoResp struct {
	domain.AuthToken
}

type ValidateEmailTokenResp struct {
	Email        string `json:"email"`
	OAuthPvodier string `json:"oauth_provider"`
}

type AuthenByEmailPassword struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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
	ActionToken  string `json:"action_token"` // New Token required consequence action
}

type AuthenByEmailPasswordDto struct {
	Email       string `json:"email"`
	RawPassword string `json:"password"`
}

type AuthenByEmailAfterSetPasswordDto struct {
	Email       string `json:"email"`
	RawPassword string `json:"password"`
	ActionToken string `json:"action_token"`
}

type ActivateUserDto struct {
	UserPkID int64 `json:"user_pkid"`
}

type AuthenByGoogleDto struct {
	Token string `json:"token"`
}

type AuthenByGoogleResponse struct {
	Profile          *domain.User `json:"profile"`
	domain.AuthToken `json:"tokens"`
}

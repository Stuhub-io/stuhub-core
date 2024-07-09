package request

type RegisterByEmailBody struct {
	Email string `json:"email" binding:"required,email"`
}

type ValidateEmailTokenBody struct {
	Token string `json:"token" binding:"required"`
}

type SetUserPasswordBody struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	ActionToken string `json:"action_token" binding:"required"`
}

type AuthenUserByEmailPasswordBody struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthenUserByGoogleBody struct {
	Token string `json:"token" binding:"required"`
}

type GetUserByTokenQuery struct {
	AccessToken string `json:"access_token" binding:"required"`
}

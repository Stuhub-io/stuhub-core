package request

type RegisterByEmailBody struct {
	Email string `binding:"required,email" json:"email"`
}

type ValidateEmailTokenBody struct {
	Token string `binding:"required" json:"token"`
}

type SetUserPasswordBody struct {
	Email       string `binding:"required,email" json:"email"`
	Password    string `binding:"required"       json:"password"`
	ActionToken string `binding:"required"       json:"action_token"`
}

type AuthenUserByEmailPasswordBody struct {
	Email    string `binding:"required,email" json:"email"`
	Password string `binding:"required"       json:"password"`
}

type AuthenUserByGoogleBody struct {
	Token string `binding:"required" json:"token"`
}

type GetUserByTokenQuery struct {
	AccessToken string `binding:"required" json:"access_token"`
}

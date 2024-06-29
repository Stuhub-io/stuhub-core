package request

type RegisterByEmailBody struct {
	Email string `json:"email" binding:"required,email"`
}

type ValidateEmailTokenBody struct {
	Token string `json:"token" binding:"required"`
}

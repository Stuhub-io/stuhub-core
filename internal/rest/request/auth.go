package request

type RegisterByEmailBody struct {
	Email string `json:"email" binding:"required,email"`
}

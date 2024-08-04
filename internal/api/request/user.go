package request

type GetUserByEmail struct {
	Email string `json:"email" binding:"required,email"`
}

type UpdateUserInfoBody struct {
	LastName  string `json:"last_name" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	Avatar    string `json:"avatar"`
}

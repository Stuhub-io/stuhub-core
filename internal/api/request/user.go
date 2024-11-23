package request

type GetUserByEmail struct {
	Email string `binding:"required,email" json:"email"`
}

type UpdateUserInfoBody struct {
	LastName  string `binding:"required" json:"last_name"`
	FirstName string `binding:"required" json:"first_name"`
	Avatar    string `json:"avatar"`
}

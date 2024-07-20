package request

type UpdateUserInfoBody struct {
	LastName  string `json:"last_name" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	Avatar    string `json:"avatar"`
}

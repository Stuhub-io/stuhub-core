package request

type GetUserByEmail struct {
	Email string `binding:"required,email" json:"email"`
}

type UpdateUserInfoBody struct {
	LastName  string `binding:"required" json:"last_name"`
	FirstName string `binding:"required" json:"first_name"`
	Avatar    string `json:"avatar"`
}

type SearchUsersBody struct {
	Search  string   `json:"search,omitempty"`
	OrgPkID *int64   `json:"org_pkid,omitempty"`
	Emails  []string `json:"emails,omitempty"`
	PaginationRequest
}

package domain

type OrganizationMember struct {
	PkId             int64  `json:"pk_id"`
	OrganizationPkID int64  `json:"organization_pkid"`
	UserPkID         int64  `json:"user_pkid"`
	Role             string `json:"role"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	// Nullable depend on usecase
	User *User
}

type Organization struct {
	PkId        int64  `json:"pk_id"`
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Role        string `json:"role"`
	Members     []OrganizationMember
}

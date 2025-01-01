package domain

type User struct {
	PkID      int64  `json:"pkid"`
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	Email     string `json:"email"`

	// Hide
	HavePassword bool   `json:"-"`
	Salt         string `json:"-"`

	// Socials
	OauthGmail string `json:"oauth_gmail"`

	ActivatedAt string `json:"activated_at"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UserSearchQuery struct {
	Limit            int      `json:"limit"`
	Offset           int      `json:"offset"`
	Search           string   `json:"search"`
	Emails           []string `json:"emails"`
	OrganizationPkID *int64   `json:"organization_pkid"`
}

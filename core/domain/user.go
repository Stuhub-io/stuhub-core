package domain

type User struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	Email     string `json:"email"`
	Password  string `json:"password"`

	// Socials
	OauthGmail bool `json:"oauth_gmail"`

	ActivatedAt string `json:"activated_at"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

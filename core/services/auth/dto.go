package auth

type AuthenByEmailDto struct {
	Email string `json:"email"`
}

type AuthenByEmailPassword struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

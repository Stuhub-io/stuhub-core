package auth

type AuthenByEmailStepOneDto struct {
	Email string `json:"email"`
}

type AuthenByEmailStepOneResp struct {
	Email           string `json:"email"`
	IsRequiredEmail bool   `json:"is_required_email"`
}

type AuthenByEmailPassword struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

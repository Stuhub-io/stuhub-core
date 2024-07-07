package domain

type OAuthProvider struct {
	Name string `json:"name"`
}

var GoogleAuthProvider = &OAuthProvider{
	Name: "google",
}

type AuthToken struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type GoogleUserInfo struct {
	Email     string
	FirstName string
	LastName  string
	Avatar    string
}

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

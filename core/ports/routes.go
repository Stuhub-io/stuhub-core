package ports

type RemoteRoute struct {
	ValidateEmailOauth    string
	ValidateOrgInvitation func(slug string) string
}

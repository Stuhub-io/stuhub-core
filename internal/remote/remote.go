package remote

import "github.com/Stuhub-io/core/ports"

func NewRemoteRoute() ports.RemoteRoute {
	return ports.RemoteRoute{
		ValidateEmailOauth: "/auth-email",
	}
}

package remote

import (
	"fmt"

	"github.com/Stuhub-io/core/ports"
)

func NewRemoteRoute() ports.RemoteRoute {
	return ports.RemoteRoute{
		ValidateEmailOauth: "/auth-email",
		ValidateOrgInvitation: func(slug string) string {
			return fmt.Sprintf("/%s/invite", slug)
		},
	}
}

package organization

import "github.com/Stuhub-io/core/domain"

type CreateOrganizationParams struct {
	OwnerPkID   int64  `json:"owner_pkid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
}

type CreateOrganizationResponse struct {
	Org *domain.Organization `json:"org"`
}

type GetRecentVisitedOrganizationParams struct {
	UserPkID int64 `json:"user_pkid"`
}

package organization

type CreateOrganizationParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Avatar      string `json:"avatar"`
	OwnerPkID   int64  `json:"owner_pkid"`
}

type GetRecentVisitedOrganizationParams struct {
	UserPkID int64 `json:"user_pkid"`
}

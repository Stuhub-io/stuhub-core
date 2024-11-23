package request

type GetSpaceByOrgPkIDParams struct {
	OrgPkID int64 `binding:"required" form:"organization_pkid"`
}

type CreateSpaceBody struct {
	OrgPkID     int64  `binding:"required" json:"organization_pkid"`
	Name        string `binding:"required" json:"name"`
	Description string `binding:"required" json:"description"`
}

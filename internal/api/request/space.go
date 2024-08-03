package request

type GetSpaceByOrgPkIDParams struct {
	OrgPkID int64 `form:"organization_pkid" binding:"required"`
}

type CreateSpaceBody struct {
	OrgPkID     int64  `json:"organization_pkid" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

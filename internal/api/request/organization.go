package request

type CreateOrgBody struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Avatar      string `json:"avatar" binding:"required"`
}

type GetOrgBySlugParams struct {
	Slug string `json:"slug" binding:"required"`
}

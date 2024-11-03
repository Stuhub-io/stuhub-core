package request

type CreatePageBody struct {
	SpacePkID      int64  `json:"space_pkid" binding:"required"`
	ViewType       string `json:"view_type" binding:"required"`
	Name           string `json:"name,omitempty"`
	ParentPagePkID *int64 `json:"parent_page_pkid,omitempty"`
	NodeID         string `json:"node_id,omitempty"`
	CoverImage     string `json:"cover_image,omitempty"`
}

type GetPagesBySpacePkIDParams struct {
	SpacePkID int64 `form:"space_pkid" binding:"required"`
}

type UpdatePageBody struct {
	Name           string `json:"name"`
	ViewType       string `json:"view_type" binding:"required"`
	ParentPagePkID *int64 `json:"parent_page_pkid,omitempty"`
	CoverImage     string `json:"cover_image,omitempty"`
	NodeID         string `json:"node_id,omitempty"`
}

type BulkGetOrCreateByNodeIDBody struct {
	PageInputs []CreatePageBody `json:"page_inputs" binding:"required"`
}

type BulkArchivePagesBody struct {
	PagePkIDs []int64 `json:"page_pkids" binding:"required"`
}

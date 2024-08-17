package request

type CreatePageBody struct {
	SpacePkID      int64  `json:"space_pk_id" binding:"required"`
	ViewType       string `json:"view_type" binding:"required"`
	Name           string `json:"name" binding:"required"`
	ParentPagePkID *int64 `json:"parent_page_pk_id,omitempty"`
}

type GetPagesBySpacePkIDParams struct {
	SpacePkID int64 `form:"space_pk_id" binding:"required"`
}

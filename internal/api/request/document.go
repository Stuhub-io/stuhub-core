package request

import "github.com/Stuhub-io/core/domain"

// page
type CreatePageBody struct {
	OrgPkID        int64               `json:"org_pkid" binding:"required"`
	ViewType       domain.PageViewType `json:"view_type" binding:"required"`
	Name           string              `json:"name,omitempty"`
	ParentPagePkID *int64              `json:"parent_page_pkid,omitempty"`
	CoverImage     string              `json:"cover_image,omitempty"`
	Document       struct {
		JsonContent string `json:"json_content,omitempty"`
	} `json:"document,omitempty"`
}

type GetPagesQuery struct {
	OrgPkID        int64                 `json:"org_pkid" form:"org_pkid" binding:"required" `
	ViewTypes      []domain.PageViewType `json:"view_types,omitempty" form:"view_types,omitempty"`
	ParentPagePkID *int64                `json:"parent_page_pkid,omitempty" form:"parent_page_pkid,omitempty"`
	IsArchived     *bool                 `json:"is_archived,omitempty" form:"is_archived,omitempty"`
	PaginationRequest
}

type UpdatePageBody struct {
	OrgPkID        *int64               `json:"org_pkid,omitempty"`
	ViewType       *domain.PageViewType `json:"view_type,omitempty"`
	Name           *string              `json:"name,omitempty"`
	ParentPagePkID *int64               `json:"parent_page_pkid,omitempty"`
	CoverImage     *string              `json:"cover_image,omitempty"`
	Document       *struct {
		JsonContent string `json:"json_content"`
	} `json:"document,omitempty"`
}

package request

import "github.com/Stuhub-io/core/domain"

// page.
type CreatePageBody struct {
	OrgPkID        int64               `binding:"required"                json:"org_pkid"`
	ViewType       domain.PageViewType `binding:"required"                json:"view_type"`
	Name           string              `json:"name,omitempty"`
	ParentPagePkID *int64              `json:"parent_page_pkid,omitempty"`
	CoverImage     string              `json:"cover_image,omitempty"`
}

type CreateDocumentBody struct {
	CreatePageBody
	Document struct {
		JsonContent string `json:"json_content,omitempty"`
	} `json:"document,omitempty"`
}

type GetPagesQuery struct {
	OrgPkID        int64                 `binding:"required"                form:"org_pkid"                   json:"org_pkid"`
	ViewTypes      []domain.PageViewType `form:"view_types,omitempty"       json:"view_types,omitempty"`
	ParentPagePkID *int64                `form:"parent_page_pkid,omitempty" json:"parent_page_pkid,omitempty"`
	IsArchived     *bool                 `form:"is_archived,omitempty"      json:"is_archived,omitempty"`
	All            bool                  `form:"all,omitempty"              json:"all,omitempty"`
	GeneralRole    *domain.PageRole      `form:"general_role,omitempty"     json:"general_role,omitempty"`
	PaginationRequest
}

type UpdatePageBody struct {
	OrgPkID    *int64               `json:"org_pkid,omitempty"`
	ViewType   *domain.PageViewType `json:"view_type,omitempty"`
	Name       *string              `json:"name,omitempty"`
	CoverImage *string              `json:"cover_image,omitempty"`
	Document   *struct {
		JsonContent string `json:"json_content"`
	} `json:"document,omitempty"`
}

type MovePageBody struct {
	ParentPagePkID *int64 `json:"parent_page_pkid,omitempty"`
}

type UpdatePageContent struct {
	JsonContent string `binding:"required" json:"json_content"`
}

type CreateAssetBody struct {
	CreatePageBody
	Asset struct {
		Url        string                `json:"url"`
		Size       int64                 `json:"size"`
		Extension  string                `json:"extension"`
		Thumbnails domain.AssetThumbnail `json:"thumbnails"`
	} `binding:"required" json:"asset"`
}

type UpdatePageGeneralAccessBody struct {
	GeneralRole domain.PageRole `binding:"required" json:"general_role,omitempty"`
}

type AddPageRoleUserBody struct {
	Role  domain.PageRole `binding:"required" json:"role,omitempty"`
	Email string          `binding:"required" json:"email"`
}

type UpdatePageRoleUserBody struct {
	Email string          `binding:"required" json:"email"`
	Role  domain.PageRole `binding:"required" json:"role,omitempty"`
}

type DeletePageRoleUserBody struct {
	Email string `binding:"required" json:"email"`
}

package activityutils

import (
	"time"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
)

type UserCreatePageMeta struct {
	ParentPagePkID *int64  `json:"parent_page_pkid"`
	ParentPageName *string `json:"parent_page_name"`
	NewPageName    string  `json:"page_name"`
	NewPagePkID    int64   `json:"page_pkid"`
	NewPageID      string  `json:"page_id"`
}

type UserMovePageMeta struct {
	OldParentPagePkID *int64  `json:"old_parent_page_pkid"`
	NewParentPagePkID *int64  `json:"new_parent_page_pkid"`
	OldParentPageName *string `json:"old_parent_page_name"`
	NewParentPageName *string `json:"new_parent_page_name"`
}

type UserRenamePageMeta struct {
	OldPageName string `json:"old_page_name"`
	NewPageName string `json:"new_page_name"`
}

type UserVisitePageMeta struct {
	ParentPagePkID int64  `json:"parent_page_pkid"`
	ParentPageID   string `json:"parent_page_id"`
	ParentPageName string `json:"parent_page_name"`
}

type UserUpdatePageInfoMeta struct {
	OldPageName  string `json:"old_page_name"`
	OldPageCover string `json:"old_page_cover"`
	OldViewType  string `json:"old_view_type"`
}

type UserRemovePageMeta struct {
	OldParentPagePkID *int64  `json:"parent_page_pkid"`
	OldParentPageName *string `json:"parent_page_name"`
}

// Meta for ActivityV2
type UserCreateFolderMeta struct {
	ParentPage *domain.Page          `json:"parent_page"`
	ChildPage  domain.Page           `json:"child_page"`
	PageRoles  []domain.PageRoleUser `json:"page_roles"`
}

type UserUploadedAssetsMeta struct {
	ParentPage *domain.Page  `json:"parent_page"`
	Assets     []domain.Page `json:"assets"`
}

type UserCreateDocumentMeta struct {
	ParentPage *domain.Page `json:"parent_page"`
	ChildPage  *domain.Page `json:"child_page"`
}

type ActivityV2ModelToDomainParams struct {
	ActivityModel    *model.Activity
	RelatedPagePkIDs []int64
}

func TransformActivityV2ModelToDomain(params ActivityV2ModelToDomainParams) *domain.ActivityV2 {
	m := params.ActivityModel
	if m == nil {
		return nil
	}

	return &domain.ActivityV2{
		PkID:             m.Pkid,
		UserPkID:         m.UserPkid,
		ActionCode:       domain.ActionCodeFromString(m.ActionCode),
		CreatedAt:        m.CreatedAt.Format(time.RFC3339),
		Snapshot:         m.Snapshot,
		RelatedPagePkIDs: params.RelatedPagePkIDs,
	}
}

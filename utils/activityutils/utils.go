package activityutils

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

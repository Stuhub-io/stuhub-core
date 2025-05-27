package domain

import "time"

const (
	ActionUserCreateFolder   ActionCode = "user.create.folder"
	ActionUserUploadedAssets ActionCode = "user.upload.assets"
	ActionUserCreateDocument ActionCode = "user.create.document"
	ActionUserRenamePage     ActionCode = "user.rename.page"
	ActionUserArchivePage    ActionCode = "user.archive.page"
)

type ActivityV2 struct {
	PkID             int64      `json:"pkid"`
	UserPkID         int64      `json:"user_pkid"`
	ActionCode       ActionCode `json:"action_code"`
	CreatedAt        string     `json:"created_at"`
	Snapshot         string     `json:"snapshot"`
	RelatedPagePkIDs []int64    `json:"related_page_pkids"`
	User             *User      `json:"user"`
}

type ActivityV2ListQuery struct {
	ActionCodes      []ActionCode `json:"action_codes"`
	UserPkIDs        []int64      `json:"user_pkids"`
	RelatedPagePkIDs []int64      `json:"related_page_pkids"`
	Limit            *int         `json:"limit"`
	EndTime          *time.Time   `json:"end_time"`
	ForUpdate        bool         `json:"for_update"`
}

type ActivityV2Input struct {
	ActionCode       ActionCode `json:"action_code"`
	UserPkID         int64      `json:"user_pkid"`
	RelatedPagePkIDs []int64    `json:"related_page_pkids"` // Important for query activity
	Snapshot         string     `json:"snapshot"`
	WithTransaction  bool       `json:"with_transaction"`
}

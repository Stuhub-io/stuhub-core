package domain

type PageAccessLog struct {
	PkID         int64  `json:"pkid"`
	Action       string `json:"action"`
	IsShared     bool   `json:"is_shared"`
	Page         Page
	ParentPages  []Page
	LastAccessed string `json:"last_accessed"`
}

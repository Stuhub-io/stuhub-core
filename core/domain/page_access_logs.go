package domain

type PageAccessLog struct {
	PkID         int64  `json:"pkid"`
	Action       string `json:"action"`
	Page         Page
	ParentPages  []Page
	LastAccessed string `json:"last_accessed"`
}

package domain

type PagePublicToken struct {
	PkID       int64  `json:"pkid"`
	ID         string `json:"id"`
	PagePkID   int64  `json:"page_pkid"`
	CreatedAt  string `json:"created_at"`
	ArchivedAt string `json:"archived_at"`
}

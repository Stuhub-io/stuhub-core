package domain

type Document struct {
	PkID        int64  `json:"pkid"`
	Content     string `json:"content"`
	JsonContent string `json:"json_content"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
	PagePkID    int64  `json:"page_pkid"`
}

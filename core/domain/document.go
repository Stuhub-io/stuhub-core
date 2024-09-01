package domain

type Document struct {
	PkId        int64  `json:"pk_id"`
	Content     string `json:"content"`
	JsonContent string `json:"json_content"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
	PagePkID    int64  `json:"page_pkid"`
}

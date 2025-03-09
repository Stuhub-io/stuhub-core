package domain

type IndexedPage struct {
	PkID           int64   `json:"pkid"`
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	AuthorPkID     int64   `json:"author_pkid"`
	AuthorFullName string  `json:"author_fullname"`
	SharedPKIDs    []int64 `json:"shared_pkids"`
	ViewType       string  `json:"view_type"`
	Content        string  `json:"content"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	ArchivedAt     string  `json:"archived_at"`
}

type SearchIndexedPageParams struct {
	UserPkID   int64
	Keyword    string
	ViewType   *string
	AuthorPkID *int64
}

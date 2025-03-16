package domain

type IndexedPage struct {
	PkID           int64   `json:"pkid,omitempty"`
	ID             string  `json:"id,omitempty"`
	Name           string  `json:"name,omitempty"`
	AuthorPkID     int64   `json:"author_pkid,omitempty"`
	AuthorFullName string  `json:"author_fullname,omitempty"`
	SharedPkIDs    []int64 `json:"shared_pkids,omitempty"`
	ViewType       string  `json:"view_type,omitempty"`
	Content        string  `json:"content,omitempty"`
	CreatedAt      string  `json:"created_at,omitempty"`
	UpdatedAt      string  `json:"updated_at,omitempty"`
	ArchivedAt     string  `json:"archived_at,omitempty"`
}

type SearchIndexedPageParams struct {
	UserPkID   int64
	Keyword    string
	ViewType   *string
	AuthorPkID *int64
}

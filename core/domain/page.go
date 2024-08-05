package domain

type Page struct {
	PkId           int64  `json:"pk_id"`
	ID             string `json:"id"`
	SpacePkID      int64  `json:"space_pkid"`
	Name           string `json:"name"`
	ParentPagePkID *int64 `json:"parent_page_pkid"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	ViewType       string `json:"view_type"`
}

type PageViewType int

const (
	PageViewTypeDoc PageViewType = iota + 1
	PageViewTypeTable
)

func (r PageViewType) String() string {
	return [...]string{"document", "table"}[r-1]
}

func PageViewFromString(val string) PageViewType {
	switch val {
	case "document":
		return PageViewTypeDoc
	case "table":
		return PageViewTypeTable
	default:
		return PageViewTypeDoc
	}
}

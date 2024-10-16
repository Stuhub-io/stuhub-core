package domain

type Page struct {
	PkID           int64  `json:"pkid"`
	ID             string `json:"id"`
	SpacePkID      int64  `json:"space_pkid"`
	Name           string `json:"name"`
	ParentPagePkID *int64 `json:"parent_page_pkid"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	ArchivedAt     string `json:"archived_at"`
	ViewType       string `json:"view_type"`
	CoverImage     string `json:"cover_image"`
}

type PageInput struct {
	Name           string `json:"name"`
	ParentPagePkID *int64 `json:"parent_page_pkid"`
	ViewType       string `json:"view_type"`
	CoverImage     string `json:"cover_image"`
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

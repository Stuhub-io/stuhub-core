package domain

import (
	"encoding/json"
	"errors"
)

type Page struct {
	PkID             int64        `json:"pkid"`
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	ParentPagePkID   *int64       `json:"parent_page_pkid"`
	AuthorPkID       *int64       `json:"author_pkid"`
	OrganizationPkID int64        `json:"organization_pkid"`
	CreatedAt        string       `json:"created_at"`
	UpdatedAt        string       `json:"updated_at"`
	ArchivedAt       string       `json:"archived_at"`
	ViewType         PageViewType `json:"view_type"`
	CoverImage       string       `json:"cover_image"`
	NodeID           string       `json:"node_id"`
	ChildPages       []Page       `json:"child_pages"`
	Document         *Document    `json:"document"`
	Asset            *Asset       `json:"asset"`
	Path             string       `json:"path"`
	IsGeneralAccess  bool         `json:"is_general_access"`
	GeneralRole      PageRole     `json:"general_role"`
	Author           *User        `json:"author"`
}

type PageRoleUser struct {
	PkID      int64    `json:"pkid"`
	PagePkID  int64    `json:"page_pkid"`
	User      User     `json:"user"`
	Role      PageRole `json:"role"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type PageInput struct {
	Name             string       `json:"name"`
	ParentPagePkID   *int64       `json:"parent_page_pkid"`
	ViewType         PageViewType `json:"view_type"`
	CoverImage       string       `json:"cover_image"`
	AuthorPkID       int64        `json:"author_pkid"`
	OrganizationPkID int64        `json:"organization_pkid"`
}

type PageUpdateInput struct {
	Name       *string       `json:"name"`
	ViewType   *PageViewType `json:"view_type"`
	CoverImage *string       `json:"cover_image"`
	Document   *struct {
		JsonContent string `json:"json_content"`
	} `json:"document"`
}

type PageMoveInput struct {
	ParentPagePkID *int64 `json:"parent_page_pkid"`
}
type PageListQuery struct {
	OrgPkID        int64          `json:"org_pkid"`
	ViewTypes      []PageViewType `json:"view_type"`
	ParentPagePkID *int64         `json:"parent_page_pkid"`
	IsArchived     *bool          `json:"is_archived"`
	Offset         int            `json:"offset"`
	Limit          int            `json:"limit"`
	IsAll          bool           `json:"all"`
}

type PageGeneralAccessUpdateInput struct {
	AuthorPkID      int64    `json:"author_pkid"`
	IsGeneralAccess bool     `json:"is_general_access"`
	GeneralRole     PageRole `json:"general_role"`
}

type PageRoleCreateInput struct {
	AuthorPkID int64    `json:"author_pkid"`
	PagePkID   int64    `json:"page_pkid"`
	UserPkID   int64    `json:"user_pkid"`
	Role       PageRole `json:"role"`
}

type PageRoleUpdateInput struct {
	AuthorPkID int64    `json:"author_pkid"`
	PagePkID   int64    `json:"page_pkid"`
	UserPkID   int64    `json:"user_pkid"`
	Role       PageRole `json:"role"`
}

type PageRoleDeleteInput struct {
	AuthorPkID int64 `json:"author_pkid"`
	PagePkID   int64 `json:"page_pkid"`
	UserPkID   int64 `json:"user_pkid"`
}

type PageRoleGetAllInput struct {
	AuthorPkID int64 `json:"author_pkid"`
	PagePkID   int64 `json:"page_pkid"`
}

type PageViewType int

const (
	PageViewTypeDoc PageViewType = iota + 1
	PageViewTypeFolder
	PageViewTypeAsset
)

func (r PageViewType) String() string {
	return [...]string{"document", "folder", "asset"}[r-1]
}

func (r *PageViewType) UnmarshalJSON(data []byte) error {
	var value int
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch PageViewType(value) {
	case PageViewTypeDoc, PageViewTypeFolder, PageViewTypeAsset:
		*r = PageViewType(value)
		return nil
	default:
		return errors.New("invalid view_type, must be 1(document) | 2(folder) | 3(asset)")
	}
}

func PageViewFromString(val string) PageViewType {
	switch val {
	case "document":
		return PageViewTypeDoc
	case "folder":
		return PageViewTypeFolder
	case "asset":
		return PageViewTypeAsset
	default:
		return PageViewTypeDoc
	}
}

type PageRole int

const (
	PageViewer PageRole = iota + 1
	PageEditor
)

func (r PageRole) String() string {
	return [...]string{"viewer", "editor"}[r-1]
}

func (r *PageRole) UnmarshalJSON(data []byte) error {
	var value int
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch PageRole(value) {
	case PageViewer, PageEditor:
		*r = PageRole(value)
		return nil
	default:
		return errors.New("invalid view_type, must be 1(viewer) | 2(editor)")
	}
}

func PageRoleFromString(val string) PageRole {
	switch val {
	case "viewer":
		return PageViewer
	case "editor":
		return PageEditor
	default:
		return PageViewer
	}
}

func (p *Page) IsAuthor(authorPkID int64) bool {
	if p.AuthorPkID == nil {
		return false
	}
	return *p.AuthorPkID == authorPkID
}

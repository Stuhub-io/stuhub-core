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
	OrganizationPkID int64        `json:"organization_pkid"`
	CreatedAt        string       `json:"created_at"`
	UpdatedAt        string       `json:"updated_at"`
	ArchivedAt       string       `json:"archived_at"`
	ViewType         PageViewType `json:"view_type"`
	CoverImage       string       `json:"cover_image"`
	NodeID           string       `json:"node_id"`
	ChildPages       []Page       `json:"child_pages"`
	Document         *Document    `json:"document"`
	Path             string       `json:"path"`
}

type Document struct {
	PkID        int64  `json:"pkid"`
	Content     string `json:"content"`
	JsonContent string `json:"json_content"`
	UpdatedAt   string `json:"updated_at"`
	CreatedAt   string `json:"created_at"`
	PagePkID    int64  `json:"page_pkid"`
}

type PageInput struct {
	Name             string        `json:"name"`
	ParentPagePkID   *int64        `json:"parent_page_pkid"`
	ViewType         PageViewType  `json:"view_type"`
	CoverImage       string        `json:"cover_image"`
	OrganizationPkID int64         `json:"organization_pkid"`
	Document         DocumentInput `json:"document"`
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

type DocumentInput struct {
	JsonContent string `json:"json_content"`
}

type PageListQuery struct {
	OrgPkID        int64          `json:"org_pkid"`
	ViewTypes      []PageViewType `json:"view_type"`
	ParentPagePkID *int64         `json:"parent_page_pkid"`
	IsArchived     *bool          `json:"is_archived"`
	Offset         int            `json:"offset"`
	Limit          int            `json:"limit"`
}

type PageViewType int

const (
	PageViewTypeDoc PageViewType = iota + 1
	PageViewTypeFolder
)

func (r PageViewType) String() string {
	return [...]string{"document", "folder"}[r-1]
}

func (r *PageViewType) UnmarshalJSON(data []byte) error {
	var value int
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch PageViewType(value) {
	case PageViewTypeDoc, PageViewTypeFolder:
		*r = PageViewType(value)
		return nil

	default:
		return errors.New("invalid view_type, must be 1(document) or 2(folder)")
	}
}

func PageViewFromString(val string) PageViewType {
	switch val {
	case "document":
		return PageViewTypeDoc
	case "folder":
		return PageViewTypeFolder
	default:
		return PageViewTypeDoc
	}
}

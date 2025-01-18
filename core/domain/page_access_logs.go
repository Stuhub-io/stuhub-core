package domain

import (
	"encoding/json"
	"errors"
	"time"
)

type PageAccessLog struct {
	PkID         int64     `json:"pkid"`
	Action       string    `json:"action"`
	IsShared     bool      `json:"is_shared"`
	Page         Page      `json:"page"`
	ParentPages  []Page    `json:"parent_pages"`
	LastAccessed time.Time `json:"last_accessed"`
}

type PageAccessAction int

const (
	PageOpen PageAccessAction = iota + 1
	PageEdit
	PageUpload
)

func (r PageAccessAction) String() string {
	return [...]string{"open", "edit", "upload"}[r-1]
}

func (r *PageAccessAction) UnmarshalJSON(data []byte) error {
	var value int
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch PageAccessAction(value) {
	case PageOpen, PageEdit, PageUpload:
		*r = PageAccessAction(value)
		return nil
	default:
		return errors.New("invalid page access action, must be 1(open) | 2(edit) | 3(upload)")
	}
}

func PageAccessActionFromString(val string) PageAccessAction {
	switch val {
	case "open":
		return PageOpen
	case "edit":
		return PageEdit
	case "upload":
		return PageUpload
	default:
		return PageOpen
	}
}

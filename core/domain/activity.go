package domain

import "time"

type ActionCode string

const (
	ActionUserCreatePage     ActionCode = "user.create.page"
	ActionUserRemovePage     ActionCode = "user.remove.page"
	ActionUserMovePage       ActionCode = "user.move.page"
	ActionUserVisitPage      ActionCode = "user.visit.page"
	ActionUserUpdatePageInfo ActionCode = "user.update.page"
)

func (a ActionCode) String() string {
	return string(a)
}

func ActionCodeFromString(s string) ActionCode {
	return ActionCode(s)
}

type Activity struct {
	ActorPkID    int64         `json:"actor_pkid"`
	Actor        *User         `json:"actor"`
	PagePkID     *int64        `json:"page_pkid"`
	Page         *Page         `json:"page"`
	OrgPkID      *int64        `json:"org_pkid"`
	Organization *Organization `json:"organization"`
	ActionCode   ActionCode    `json:"action_code"`
	Label        *string       `json:"label"`
	CreatedAt    string        `json:"created_at"`
	UpdatedAt    string        `json:"updated_at"`
	MetaData     *string       `json:"meta_data"`
}

type ActivityListQuery struct {
	ActionCodes   []ActionCode `json:"action_codes"`
	ActorPkIDs    []int64      `json:"actor_pkids"`
	PagePkIDs     []int64      `json:"page_pkids"`
	PageToken     string       `json:"page_token,omitempty"`     // Token for continuation from previous page
	SortDirection string       `json:"sort_direction,omitempty"` // "asc" or "desc", default is "desc" for most recent first
	StartTime     *time.Time   `json:"start_time"`
	EndTime       *time.Time   `json:"end_time"`
	Limit         int          `json:"limit"`
}

type ActivityInput struct {
	ActionCode ActionCode `json:"action_code"`
	ActorPkID  int64      `json:"actor_pkid"`
	PagePkID   *int64     `json:"page_pkid"`
	OrgPkID    *int64     `json:"org_pkid"`
	Label      *string    `json:"label"`
	MetaData   *string    `json:"meta_data"`
}

type ActivityMetaParams struct {
	ActionCode ActionCode
	Actor      *User
	Page       *Page
}

// Query All child page activities
type PageActivitiesListQuery struct {
	ActionCodes []ActionCode `json:"action_codes"`
	PagePkID    int64
}

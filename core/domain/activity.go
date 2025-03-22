package domain

type ActionCode string

const (
	ActionUserCreatePage           ActionCode = "user.create.page"
	ActionUserModifiedPageDocument ActionCode = "user.modify.page.document"
	ActionUserRemovePage           ActionCode = "user.remove.page"
	ActionUserVisitPage            ActionCode = "user.visit.page"
)

func (a ActionCode) String() string {
	return string(a)
}
func ActionCodeFromString(s string) ActionCode {
	return ActionCode(s)
}

type Activity struct {
	PkID       int64      `json:"pkid"`
	ActorPkID  int64      `json:"actor_pkid"`
	Actor      *User      `json:"actor"`
	PagePkID   *int64     `json:"page_pkid"`
	Page       *Page      `json:"page"`
	ActionCode ActionCode `json:"action_code"`
	Label      *string    `json:"label"`
	CreatedAt  string     `json:"created_at"`
	UpdatedAt  string     `json:"updated_at"`
	MetaData   *string    `json:"meta_data"`
}

type ActivityListQuery struct {
	PkIDs       []int64      `json:"pkids"`
	ActionCodes []ActionCode `json:"action_code"`
	ActorPkIDs  []int64      `json:"actor_pkid"`
	PagePkIDs   []int64      `json:"page_pkid"`
}

type ActivityInput struct {
	ActionCode ActionCode `json:"action_code"`
	ActorPkID  int64      `json:"actor_pkid"`
	PagePkID   *int64     `json:"page_pkid"`
	Label      *string    `json:"label"`
	MetaData   *string    `json:"meta_data"`
}

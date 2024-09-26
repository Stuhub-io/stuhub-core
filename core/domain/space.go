package domain

type Space struct {
	PkID        int64         `json:"pkid"`
	ID          string        `json:"id"`
	OrgPkID     int64         `json:"org_pkid"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	IsPrivate   bool          `json:"is_private"`
	CreatedAt   string        `json:"created_at"`
	UpdatedAt   string        `json:"updated_at"`
	Members     []SpaceMember `json:"members"`
}

type SpaceMemberRole int

const (
	RoleSpaceOwner SpaceMemberRole = iota + 1
	RoleSpaceMember
	RoleSpaceGuest
)

func (r SpaceMemberRole) String() string {
	return [...]string{"owner", "member", "guest"}[r-1]
}

type SpaceMember struct {
	PkID      int64  `json:"pkid"`
	SpacePkID int64  `json:"space_pkid"`
	UserPkID  int64  `json:"user_pkid"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	User      *User  `json:"user"`
}

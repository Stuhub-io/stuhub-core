package domain

type Organization struct {
	PkID        int64                `json:"pkid"`
	ID          string               `json:"id"`
	OwnerID     int64                `json:"owner_id"`
	Name        string               `json:"name"`
	Slug        string               `json:"slug"`
	Description string               `json:"description"`
	Avatar      string               `json:"avatar"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
	Members     []OrganizationMember `json:"members"`
}

type OrganizationMemberRole int

const (
	Owner OrganizationMemberRole = iota + 1
	Member
)

func (r OrganizationMemberRole) String() string {
	return [...]string{"owner", "member"}[r-1]
}

type OrganizationMember struct {
	PkID             int64  `json:"pkid"`
	OrganizationPkID int64  `json:"organization_pkid"`
	UserPkID         *int64 `json:"user_pkid"`
	Role             string `json:"role"`
	ActivatedAt      string `json:"activated_at"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	User             *User  `json:"user"`
}

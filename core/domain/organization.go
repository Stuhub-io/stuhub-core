package domain

type Organization struct {
	PkId        int64                `json:"-"`
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
	PkId             int64  `json:"pk_id"`
	OrganizationPkID int64  `json:"organization_pkid"`
	UserPkID         *int64 `json:"user_pkid"`
	Role             string `json:"role"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	// Nullable depend on usecase
	User *User `json:"user"`
}

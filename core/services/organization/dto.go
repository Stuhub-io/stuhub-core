package organization

import "github.com/Stuhub-io/core/domain"

type CreateOrganizationDto struct {
	OwnerPkID   int64
	Name        string
	Description string
	Avatar      string
}

type CreateOrganizationResponse struct {
	Org *domain.Organization `json:"org"`
}

type GetRecentVisitedOrganizationDto struct {
	UserPkID int64
}

type OrgInviteInfo struct {
	PkID    int64  `binding:"required" json:"pkid"`
	Name    string `binding:"required" json:"name"`
	Slug    string `binding:"required" json:"slug"`
	Members int64  `binding:"required" json:"members"`
	Avatar  string `binding:"required" json:"avatar"`
}

type EmailInviteInfo struct {
	Email string `binding:"required" json:"email"`
	Role  string `binding:"required" json:"role"`
}

type InviteMemberByEmailsDto struct {
	Owner       *domain.User
	OrgInfo     OrgInviteInfo
	InviteInfos []EmailInviteInfo
}

type InviteMemberByEmailsResponse struct {
	SentEmails   []string `json:"sent_emails"`
	FailedEmails []string `json:"failed_emails"`
}

type ValidateOrgInviteTokenDto struct {
	UserPkID int64
	Token    string
}

type AddMemberToOrgDto struct {
	UserPkID int64
	OrgPkID  int64
	Role     domain.OrganizationMemberRole
}

type ActivateMemberDto struct {
	MemberPkID int64
	OrgPkID    int64
}

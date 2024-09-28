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
	PkID    int64  `json:"pkid" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Slug    string `json:"slug" binding:"required"`
	Members int64  `json:"members" binding:"required"`
	Avatar  string `json:"avatar" binding:"required"`
}

type EmailInviteInfo struct {
	Email string `json:"email" binding:"required"`
	Role  string `json:"role" binding:"required"`
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

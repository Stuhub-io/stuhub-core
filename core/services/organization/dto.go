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

type EmailInviteInfo struct {
	Email string
	Role  string
}

type InviteMemberByEmailsDto struct {
	Emails []EmailInviteInfo
}

type AddMemberToOrgDto struct {
	UserPkID int64
	OrgPkID  int64
	Role     domain.OrganizationMemberRole
}

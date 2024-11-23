package request

import "github.com/Stuhub-io/core/services/organization"

type CreateOrgBody struct {
	Name        string `binding:"required" json:"name"`
	Description string `binding:"required" json:"description"`
	Avatar      string `binding:"required" json:"avatar"`
}

type GetOrgBySlugParams struct {
	Slug string `binding:"required" form:"slug"`
}

type InviteMembersByEmailParams struct {
	OrgInfo organization.OrgInviteInfo     `binding:"required" json:"org_info"`
	Infos   []organization.EmailInviteInfo `binding:"required" json:"infos"`
}

type ValidateOrgInvitationParams struct {
	Token string `binding:"required" json:"token"`
}

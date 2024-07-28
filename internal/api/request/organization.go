package request

import "github.com/Stuhub-io/core/services/organization"

type CreateOrgBody struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	Avatar      string `json:"avatar" binding:"required"`
}

type GetOrgBySlugParams struct {
	Slug string `json:"slug" binding:"required"`
}

type InviteMembersByEmailParams struct {
	OrgInfo organization.OrgInviteInfo     `json:"org_info" binding:"required"`
	Infos   []organization.EmailInviteInfo `json:"infos" binding:"required"`
}

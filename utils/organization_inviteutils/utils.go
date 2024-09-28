package organization_inviteutils

import (
	"strconv"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/organizationutils"
	"github.com/gin-gonic/gin"
)

const (
	OrganizationInvitePkIDParam = "organizationInvitePkIDParam"
)

const InviteIDParam = "inviteID"

func GetInviteIDParam(c *gin.Context) (string, bool) {
	inviteID := c.Params.ByName(InviteIDParam)
	if inviteID == "" {
		return "", false
	}
	return inviteID, true
}

type InviteWithOrganization struct {
	model.OrganizationInvite
	Organization organizationutils.OrganizationWithMembers `gorm:"foreignKey:organization_pkid" json:"organization"` // Define foreign key relationship
}

func GetOrganizationInviteParams(c *gin.Context) (int64, bool) {
	param := c.Params.ByName(OrganizationInvitePkIDParam)
	if param == "" {
		return int64(-1), false
	}

	organizationInvitePkID, cErr := strconv.Atoi(param)

	return int64(organizationInvitePkID), cErr == nil
}

func TransformOrganizationInviteModelToDomain(invite model.OrganizationInvite) *domain.OrganizationInvite {
	return &domain.OrganizationInvite{
		PkID:             invite.Pkid,
		ID:               invite.ID,
		UserPkID:         invite.UserPkid,
		OrganizationPkID: invite.OrganizationPkid,
		IsUsed:           invite.IsUsed,
		CreatedAt:        invite.CreatedAt,
		ExpiredAt:        invite.ExpiredAt,
	}
}

func TransformOrganizationInviteModelToDomain_WithOrg(invite InviteWithOrganization) *domain.OrganizationInvite {
	return &domain.OrganizationInvite{
		PkID:             invite.Pkid,
		ID:               invite.ID,
		UserPkID:         invite.UserPkid,
		OrganizationPkID: invite.OrganizationPkid,
		IsUsed:           invite.IsUsed,
		Organization:     organizationutils.TransformOrganizationModelToDomain(invite.Organization),
		CreatedAt:        invite.CreatedAt,
		ExpiredAt:        invite.ExpiredAt,
	}
}

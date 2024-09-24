package organzation_inviteutils

import (
	"strconv"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/gin-gonic/gin"
)

const (
	OrganizationInvitePkIDParam = "organizationInvitePkIDParam"
)

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
		PkId:             invite.Pkid,
		ID:               invite.ID,
		UserPkID:         invite.UserPkid,
		OrganizationPkID: invite.OrganizationPkid,
		IsUsed:           invite.IsUsed,
		CreatedAt:        invite.CreatedAt,
		ExpiredAt:        invite.ExpiredAt,
	}
}

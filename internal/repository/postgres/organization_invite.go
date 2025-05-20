package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	organization_inviteutils "github.com/Stuhub-io/utils/organization_inviteutils"
	"gorm.io/gorm"
)

type OrganizationInvitesRepository struct {
	DB *DB
}

type NewOrganizationInvitesRepositoryParams struct {
	DB *DB
}

func NewOrganizationInvitesRepository(DB *DB) *OrganizationInvitesRepository {
	return &OrganizationInvitesRepository{
		DB: DB,
	}
}

func (r *OrganizationInvitesRepository) CreateInvite(ctx context.Context, organizationPkID int64, userPkID int64) (*domain.OrganizationInvite, *domain.Error) {
	newInvite := model.OrganizationInvite{
		UserPkid:         userPkID,
		OrganizationPkid: organizationPkID,
		ExpiredAt:        time.Now().Add(domain.OrgInvitationExpiredTime),
	}

	err := r.DB.DB().Create(&newInvite).Error
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	return organization_inviteutils.TransformOrganizationInviteModelToDomain(newInvite), nil
}

func (r *OrganizationInvitesRepository) UpdateInvite(ctx context.Context, invite model.OrganizationInvite) (*domain.OrganizationInvite, *domain.Error) {
	var updatedInvite model.OrganizationInvite

	err := r.DB.DB().Model(&updatedInvite).Where("id = ?", invite.ID).Update("is_used", true).Error

	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	return organization_inviteutils.TransformOrganizationInviteModelToDomain(updatedInvite), nil
}

func (r *OrganizationInvitesRepository) GetInviteByID(ctx context.Context, inviteID string) (*domain.OrganizationInvite, *domain.Error) {
	var invite organization_inviteutils.InviteWithOrganization

	err := r.DB.DB().Preload("Organization.Members").Preload("Organization.Owner").Where("id = ?", inviteID).First(&invite).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrOrgNotFound
		}
		return nil, domain.ErrDatabaseQuery
	}

	return organization_inviteutils.TransformOrganizationInviteModelToDomain_WithOrg(invite), nil
}

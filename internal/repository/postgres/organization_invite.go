package postgres

import (
	"context"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/model"
	organization_inviteutils "github.com/Stuhub-io/utils/organization_inviteutils"
)

type OrganizationInvitesRepository struct {
	store *store.DBStore
	cfg   config.Config
}

type NewOrganizationInvitesRepositoryParams struct {
	Cfg   config.Config
	Store *store.DBStore
}

func NewOrganizationInvitesRepository(params NewOrganizationInvitesRepositoryParams) *OrganizationInvitesRepository {
	return &OrganizationInvitesRepository{
		store: params.Store,
		cfg:   params.Cfg,
	}
}

func (r *OrganizationInvitesRepository) CreateInvite(ctx context.Context, organizationPkID int64, userPkID int64) (*domain.OrganizationInvite, *domain.Error) {
	newInvite := model.OrganizationInvite{
		UserPkid:         userPkID,
		OrganizationPkid: organizationPkID,
		ExpiredAt:        time.Now().Add(domain.OrgInvitationExpiredTime),
	}

	err := r.store.DB().Create(&newInvite).Error
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	return organization_inviteutils.TransformOrganizationInviteModelToDomain(newInvite), nil
}

func (r *OrganizationInvitesRepository) UpdateInvite(ctx context.Context, invite model.OrganizationInvite) (*domain.OrganizationInvite, *domain.Error) {
	var updatedInvite model.OrganizationInvite

	err := r.store.DB().Model(&updatedInvite).Updates(invite).Error
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	return organization_inviteutils.TransformOrganizationInviteModelToDomain(updatedInvite), nil
}

func (r *OrganizationInvitesRepository) GetInviteByID(ctx context.Context, inviteID string) (*domain.OrganizationInvite, *domain.Error) {
	var invite organization_inviteutils.InviteWithOrganization

	err := r.store.DB().Preload("Organization.Members").Preload("Organization.Owner").Where("id = ?", inviteID).First(&invite).Error
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	return organization_inviteutils.TransformOrganizationInviteModelToDomain_WithOrg(invite), nil
}

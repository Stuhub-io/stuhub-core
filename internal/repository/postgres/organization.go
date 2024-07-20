package postgres

import (
	"context"
	"errors"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/model"
	commonutils "github.com/Stuhub-io/utils"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	store *store.DBStore
	cfg   config.Config
}

type NewOrganizationRepositoryParams struct {
	Store *store.DBStore
	Cfg   config.Config
}

func NewOrganizationRepository(params NewOrganizationRepositoryParams) ports.OrganizationRepository {
	return &OrganizationRepository{
		store: params.Store,
		cfg:   params.Cfg,
	}
}

func (r *OrganizationRepository) GetOrgMember(ctx context.Context, pkID int64, includeUser bool) ([]domain.OrganizationMember, *domain.Error) {
	var members []model.OrganizationMember
	err := r.store.DB().Where("organization_id = ?", pkID).Find(&members).Error
	if err != nil {
		return nil, domain.ErrInternalServerError
	}
	var orgMembesrResp = make([]domain.OrganizationMember, len(members))
	for i, member := range members {
		orgMembesrResp[i] = domain.OrganizationMember{
			PkId:             member.Pkid,
			OrganizationPkID: member.OrganizationPkid,
			UserPkID:         member.UserPkid,
			Role:             member.Role,
			CreatedAt:        member.CreatedAt.String(),
			UpdatedAt:        member.UpdatedAt.String(),
			// Include User entity if `includeUser` is true
			// User: domain.User{...},
		}
	}
	return orgMembesrResp, nil
}

func (r *OrganizationRepository) GetOrgBySlug(ctx context.Context, slug string) (*domain.Organization, *domain.Error) {
	var org model.Organization
	err := r.store.DB().Where("slug = ?", slug).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
	}

	member, err2 := r.GetOrgMember(ctx, org.Pkid, false)
	if err2 != nil {
		return nil, err2
	}
	var currentMember domain.OrganizationMember
	for _, mem := range member {
		if mem.UserPkID == 1 {
			currentMember = mem
			break
		}
	}

	return &domain.Organization{
		PkId:        org.Pkid,
		ID:          org.ID,
		Name:        org.Name,
		Slug:        org.Slug,
		Description: org.Description,
		Avatar:      org.Avatar,
		CreatedAt:   org.CreatedAt.String(),
		UpdatedAt:   org.UpdatedAt.String(),
		Role:        currentMember.Role,
		Members:     member,
	}, nil
}

func (r *OrganizationRepository) CreateOrg(ctx context.Context, ownerPkID int64, name, description, avatar string) (*domain.Organization, *domain.Error) {
	// FIXME: Need redis lock to prevent collision
	slugText := slug.Make(name)
	var existSlutOrg []model.Organization
	err := r.store.DB().Where("slug LIKE ?", slugText+"%").Find(&existSlutOrg).Error
	if err != nil {
		return nil, domain.ErrInternalServerError
	}
	existSlugs := make([]string, len(existSlutOrg))
	for i, org := range existSlutOrg {
		existSlugs[i] = org.Slug
	}
	cleanSlug := commonutils.GetSlugResolution(existSlugs, slugText)
	//

	newOrg := model.Organization{
		Name:   name,
		Avatar: avatar,
		Slug:   cleanSlug,
	}

	err = r.store.DB().Create(&newOrg).Error
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	return nil, nil
}

func (r *OrganizationRepository) GetOrgsByUserPkID(ctx context.Context, userPkID int64) ([]domain.Organization, *domain.Error) {
	return nil, nil
}

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
	cfg            config.Config
	store          *store.DBStore
	userRepository ports.UserRepository
}

type NewOrganizationRepositoryParams struct {
	Cfg            config.Config
	Store          *store.DBStore
	UserRepository ports.UserRepository
}

type MemberWithUser struct {
	model.OrganizationMember
	User model.User `gorm:"foreignKey:user_pkid" json:"user"` // Define foreign key relationship
}

type OrganizationWithMembers struct {
	model.Organization
	Members []MemberWithUser `gorm:"foreignKey:organization_pkid" json:"members"` // Consider JSON tag for future use
}

func NewOrganizationRepository(params NewOrganizationRepositoryParams) ports.OrganizationRepository {
	return &OrganizationRepository{
		cfg:            params.Cfg,
		store:          params.Store,
		userRepository: params.UserRepository,
	}
}

func mapOrg(model model.Organization, members []domain.OrganizationMember) *domain.Organization {
	return &domain.Organization{
		ID:          model.ID,
		PkId:        model.Pkid,
		OwnerID:     model.OwnerID,
		Name:        model.Name,
		Slug:        model.Slug,
		Description: model.Description,
		Avatar:      model.Avatar,
		CreatedAt:   model.CreatedAt.String(),
		UpdatedAt:   model.UpdatedAt.String(),
		Members:     members,
	}
}

func mapOrgModelToDomain(model OrganizationWithMembers, members []domain.OrganizationMember) *domain.Organization {
	return &domain.Organization{
		ID:          model.ID,
		PkId:        model.Pkid,
		OwnerID:     model.OwnerID,
		Name:        model.Name,
		Slug:        model.Slug,
		Description: model.Description,
		Avatar:      model.Avatar,
		CreatedAt:   model.CreatedAt.String(),
		UpdatedAt:   model.UpdatedAt.String(),
		Members:     members,
	}
}

func mapOrgModelsToDomain(models []OrganizationWithMembers) []*domain.Organization {
	domainOrgs := make([]*domain.Organization, 0, len(models))

	for _, org := range models {
		domainMembers := make([]domain.OrganizationMember, 0, len(org.Members))
		for _, member := range org.Members {
			domainMember := domain.OrganizationMember{
				PkId:             member.Pkid,
				OrganizationPkID: member.OrganizationPkid,
				UserPkID:         member.Pkid,
				Role:             member.Role,
				User:             mapUserModelToDomain(member.User),
				CreatedAt:        member.CreatedAt.String(),
				UpdatedAt:        member.UpdatedAt.String(),
			}

			domainMembers = append(domainMembers, domainMember)
		}

		domainOrg := mapOrgModelToDomain(org, domainMembers)
		domainOrgs = append(domainOrgs, domainOrg)
	}

	return domainOrgs
}

func mapOrgMemberModelToDomain(model model.OrganizationMember, user *domain.User) *domain.OrganizationMember {
	return &domain.OrganizationMember{
		PkId:             model.Pkid,
		OrganizationPkID: model.OrganizationPkid,
		UserPkID:         model.UserPkid,
		Role:             model.Role,
		User:             user,
		CreatedAt:        model.CreatedAt.String(),
		UpdatedAt:        model.UpdatedAt.String(),
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

func (r *OrganizationRepository) GetOwnerOrgByName(ctx context.Context, ownerID int64, name string) (*domain.Organization, *domain.Error) {
	var org model.Organization

	err := r.store.DB().Where("owner_id = ? AND name = ?", ownerID, name).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrOrgNotFound
		}

		return nil, domain.ErrDatabaseQuery
	}

	return &domain.Organization{
		PkId: org.Pkid,
	}, nil
}

func (r *OrganizationRepository) GetOrgBySlug(ctx context.Context, slug string) (*domain.Organization, *domain.Error) {
	var org model.Organization

	err := r.store.DB().Where("slug = ?", slug).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
	}

	return &domain.Organization{
		PkId: org.Pkid,
	}, nil
}

func (r *OrganizationRepository) CreateOrg(ctx context.Context, ownerPkID int64, name, description, avatar string) (*domain.Organization, *domain.Error) {
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

	var newOrg model.Organization
	var ownerMember model.OrganizationMember
	err = r.store.DB().Transaction(func(tx *gorm.DB) error {
		newOrg = model.Organization{
			OwnerID:     ownerPkID,
			Name:        name,
			Description: description,
			Avatar:      avatar,
			Slug:        cleanSlug,
		}
		err = r.store.DB().Create(&newOrg).Error
		if err != nil {
			return err
		}

		ownerMember = model.OrganizationMember{
			OrganizationPkid: newOrg.Pkid,
			UserPkid:         ownerPkID,
			Role:             domain.Owner.String(),
		}
		err = r.store.DB().Create(&ownerMember).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	owner, uerr := r.userRepository.GetUserByPkID(context.Background(), ownerMember.UserPkid)
	if uerr != nil {
		return nil, uerr
	}

	return mapOrg(newOrg, []domain.OrganizationMember{
		*mapOrgMemberModelToDomain(ownerMember, owner),
	}), nil
}

func (r *OrganizationRepository) GetOrgsByUserPkID(ctx context.Context, userPkID int64) ([]*domain.Organization, *domain.Error) {
	var joinedOrgs []OrganizationWithMembers

	err := r.store.DB().Preload("Members.User").
		Joins("JOIN organization_member ON organization_member.organization_pkid = organizations.pkid").
		Where("organization_member.user_pkid = ?", userPkID).
		Find(&joinedOrgs).Error
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	return mapOrgModelsToDomain(joinedOrgs), nil
}

func (r *OrganizationRepository) AddMemberToOrg(ctx context.Context, userPkID, orgPkID int64, role domain.OrganizationMemberRole) (*domain.OrganizationMember, *domain.Error) {
	var existingMember model.OrganizationMember

	rel := r.store.DB().Where("organization_pkid = ? AND user_pkid = ?").First(&existingMember)
	if err := rel.Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrDatabaseQuery
	}

	if rel.RowsAffected == 1 {
		return nil, domain.ErrExistOrgMember(userPkID)
	}

	var newMember = model.OrganizationMember{
		OrganizationPkid: orgPkID,
		UserPkid:         userPkID,
		Role:             role.String(),
	}
	err := r.store.DB().Create(&newMember).Error
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	user, dErr := r.userRepository.GetUserByPkID(context.Background(), userPkID)
	if dErr != nil {
		return nil, dErr
	}

	return mapOrgMemberModelToDomain(newMember, user), nil
}

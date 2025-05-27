package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/internal/repository/model"
	commonutils "github.com/Stuhub-io/utils"
	"github.com/Stuhub-io/utils/organizationutils"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrganizationRepository struct {
	cfg            config.Config
	DB             *DB
	userRepository ports.UserRepository
}

type NewOrganizationRepositoryParams struct {
	Cfg            config.Config
	DB             *DB
	UserRepository ports.UserRepository
}

func NewOrganizationRepository(DB *DB, userRepo ports.UserRepository) ports.OrganizationRepository {
	return &OrganizationRepository{
		DB:             DB,
		userRepository: userRepo,
	}
}

func (r *OrganizationRepository) GetOrgMembers(ctx context.Context, pkID int64) ([]domain.OrganizationMember, *domain.Error) {
	var members []organizationutils.MemberWithUser

	err := r.DB.DB().Preload("User").Where("organization_id = ?", pkID).First(&members).Error
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	return organizationutils.TransformOrganizationMemberModelToDomain_Many(members), nil
}

func (r *OrganizationRepository) GetOrgMemberByEmail(ctx context.Context, orgPkID int64, email string) (*domain.OrganizationMember, *domain.Error) {
	var member organizationutils.MemberWithUser

	err := r.DB.DB().Preload("User").
		Joins("JOIN users ON users.pkid = organization_member.user_pkid").
		Where("organization_pkid = ? AND users.email = ?", orgPkID, email).
		First(&member).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrOrgMemberNotFound
		}
		return nil, domain.ErrInternalServerError
	}

	return organizationutils.TransformOrganizationMemberModelToDomain(member), nil
}

func (r *OrganizationRepository) GetOrgMemberByUserPkID(ctx context.Context, orgPkID int64, userPkID int64) (*domain.OrganizationMember, *domain.Error) {
	var member organizationutils.MemberWithUser

	err := r.DB.DB().Preload("User").
		Joins("JOIN users ON users.pkid = organization_member.user_pkid").
		Where("organization_pkid = ? AND users.pkid = ?", orgPkID, userPkID).
		First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrOrgMemberNotFound
		}
		return nil, domain.ErrInternalServerError
	}

	return organizationutils.TransformOrganizationMemberModelToDomain(member), nil
}

func (r *OrganizationRepository) GetOwnerOrgByName(ctx context.Context, ownerID int64, name string) (*domain.Organization, *domain.Error) {
	var org organizationutils.OrganizationWithMembers

	err := r.DB.DB().Preload("Members").Where("owner_id = ? AND name = ?", ownerID, name).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrOrgNotFound
		}

		return nil, domain.ErrDatabaseQuery
	}

	return organizationutils.TransformOrganizationModelToDomain(org), nil
}

func (r *OrganizationRepository) GetOwnerOrgByPkID(ctx context.Context, ownerID, pkID int64) (*domain.Organization, *domain.Error) {
	var org organizationutils.OrganizationWithMembers

	err := r.DB.DB().Preload("Members").Where("owner_id = ? AND pkid = ?", ownerID, pkID).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrOrgNotFound
		}

		return nil, domain.ErrDatabaseQuery
	}

	return organizationutils.TransformOrganizationModelToDomain(org), nil
}

func (r *OrganizationRepository) GetOrgByPkID(ctx context.Context, pkID int64) (*domain.Organization, *domain.Error) {
	var org organizationutils.OrganizationWithMembers

	err := r.DB.DB().Preload("Members").Where("pkid = ?", pkID).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
	}

	return organizationutils.TransformOrganizationModelToDomain(org), nil
}

func (r *OrganizationRepository) GetOrgBySlug(ctx context.Context, slug string) (*domain.Organization, *domain.Error) {
	var org organizationutils.OrganizationWithMembers

	err := r.DB.DB().Preload("Members").Where("slug = ?", slug).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
	}

	return organizationutils.TransformOrganizationModelToDomain(org), nil
}

func (r *OrganizationRepository) CreateOrg(ctx context.Context, ownerPkID int64, name, description, avatar string) (*domain.Organization, *domain.Error) {
	slugText := slug.Make(name)

	var existSlutOrg []model.Organization

	err := r.DB.DB().Where("slug LIKE ?", slugText+"%").Find(&existSlutOrg).Error
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

	// -- Transaction
	tx, doneTx := r.DB.WithOptionalTransaction(true)
	newOrg = model.Organization{
		OwnerID:     ownerPkID,
		Name:        name,
		Description: description,
		Avatar:      avatar,
		Slug:        cleanSlug,
	}
	err = tx.DB().Create(&newOrg).Error
	if err != nil {
		return nil, doneTx(err)
	}

	ownerMember = model.OrganizationMember{
		OrganizationPkid: newOrg.Pkid,
		UserPkid:         &ownerPkID,
		Role:             domain.Owner.String(),
	}

	err = tx.DB().Create(&ownerMember).Error
	if err != nil {
		return nil, doneTx(err)
	}
	commitErr := doneTx(nil)
	if commitErr != nil {
		return nil, commitErr
	}
	// -- End Transaction

	owner, uerr := r.userRepository.GetUserByPkID(context.Background(), *ownerMember.UserPkid)
	if uerr != nil {
		return nil, uerr
	}

	return organizationutils.TransformOrganizationModelToDomain_New(newOrg, ownerMember, owner), nil
}

func (r *OrganizationRepository) GetOrgsByUserPkID(ctx context.Context, userPkID int64) ([]*domain.Organization, *domain.Error) {
	var joinedOrgs []organizationutils.OrganizationWithMembers

	err := r.DB.DB().Preload("Members.User").
		Joins("JOIN organization_member ON organization_member.organization_pkid = organizations.pkid").
		Where("organization_member.user_pkid = ?", userPkID).
		Find(&joinedOrgs).Error
	if err != nil {
		return nil, domain.ErrDatabaseQuery
	}

	return organizationutils.TransformOrganizationModelToDomain_Many(joinedOrgs), nil
}

func (r *OrganizationRepository) AddMemberToOrg(ctx context.Context, orgPkID int64, userPkID *int64, role string) (*domain.OrganizationMember, *domain.Error) {
	var newMember = model.OrganizationMember{
		OrganizationPkid: orgPkID,
		UserPkid:         userPkID,
		Role:             role,
	}
	err := r.DB.DB().Create(&newMember).Error
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	var user *domain.User

	if newMember.UserPkid != nil {
		user, _ = r.userRepository.GetUserByPkID(context.Background(), *newMember.UserPkid)
	}

	return organizationutils.TransformOrganizationMemberModelToDomain_New(newMember, user), nil
}

func (r *OrganizationRepository) SetOrgMemberActivatedAt(ctx context.Context, pkID int64, activatedAt time.Time) (*domain.OrganizationMember, *domain.Error) {
	var member model.OrganizationMember

	err := r.DB.DB().Model(&member).Clauses(clause.Returning{}).Where("pkid = ?", pkID).Update("activated_at", activatedAt).Error
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	return organizationutils.TransformOrganizationMemberModelToDomain_New(member, nil), nil
}

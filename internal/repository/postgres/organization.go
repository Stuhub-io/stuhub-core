package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/model"
	commonutils "github.com/Stuhub-io/utils"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
			activatedAt := ""
			if member.ActivatedAt != nil {
				activatedAt = member.ActivatedAt.String()
			}

			domainMember := domain.OrganizationMember{
				PkId:             member.Pkid,
				OrganizationPkID: member.OrganizationPkid,
				UserPkID:         &member.Pkid,
				Role:             member.Role,
				User:             mapUserModelToDomain(member.User),
				ActivatedAt:      activatedAt,
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

func mapOrgMemberUserModelToDomain(model model.OrganizationMember, user *domain.User) *domain.OrganizationMember {
	activatedAt := ""
	if model.ActivatedAt != nil {
		activatedAt = model.ActivatedAt.String()
	}

	return &domain.OrganizationMember{
		PkId:             model.Pkid,
		OrganizationPkID: model.OrganizationPkid,
		UserPkID:         model.UserPkid,
		Role:             model.Role,
		User:             user,
		ActivatedAt:      activatedAt,
		CreatedAt:        model.CreatedAt.String(),
		UpdatedAt:        model.UpdatedAt.String(),
	}
}

func mapOrgMemberModelsToDomain(models []MemberWithUser) []*domain.OrganizationMember {
	domainOrgMembers := make([]*domain.OrganizationMember, 0, len(models))
	for _, member := range models {
		activatedAt := ""
		if member.ActivatedAt != nil {
			activatedAt = member.ActivatedAt.String()
		}
		domainMember := &domain.OrganizationMember{
			PkId:             member.Pkid,
			OrganizationPkID: member.OrganizationPkid,
			UserPkID:         &member.Pkid,
			Role:             member.Role,
			User:             mapUserModelToDomain(member.User),
			ActivatedAt:      activatedAt,
			CreatedAt:        member.CreatedAt.String(),
			UpdatedAt:        member.UpdatedAt.String(),
		}
		domainOrgMembers = append(domainOrgMembers, domainMember)
	}

	return domainOrgMembers
}

func (r *OrganizationRepository) GetOrgMembers(ctx context.Context, pkID int64) ([]*domain.OrganizationMember, *domain.Error) {
	var members []MemberWithUser

	err := r.store.DB().Preload("User").Where("organization_id = ?", pkID).First(&members).Error
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	return mapOrgMemberModelsToDomain(members), nil
}

func (r *OrganizationRepository) GetOrgMemberByEmail(ctx context.Context, orgPkId int64, email string) (*domain.OrganizationMember, *domain.Error) {
	var member MemberWithUser

	err := r.store.DB().Preload("User").
		Joins("JOIN users ON users.pkid = organization_member.user_pkid").
		Where("organization_pkid = ? AND users.email = ?", orgPkId, email).
		First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrOrgMemberNotFound
		}
		return nil, domain.ErrInternalServerError
	}

	activatedAt := ""
	if member.ActivatedAt != nil {
		activatedAt = member.ActivatedAt.String()
	}

	return &domain.OrganizationMember{
		PkId:             member.Pkid,
		OrganizationPkID: member.OrganizationPkid,
		UserPkID:         member.UserPkid,
		Role:             member.Role,
		User:             mapUserModelToDomain(member.User),
		ActivatedAt:      activatedAt,
		CreatedAt:        member.CreatedAt.String(),
		UpdatedAt:        member.UpdatedAt.String(),
	}, nil
}

func (r *OrganizationRepository) GetOrgMemberByUserPkID(ctx context.Context, orgPkId int64, userPkId int64) (*domain.OrganizationMember, *domain.Error) {
	var member MemberWithUser

	err := r.store.DB().Preload("User").
		Joins("JOIN users ON users.pkid = organization_member.user_pkid").
		Where("organization_pkid = ? AND users.pkid = ?", orgPkId, userPkId).
		First(&member).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrOrgMemberNotFound
		}
		return nil, domain.ErrInternalServerError
	}

	activatedAt := ""
	if member.ActivatedAt != nil {
		activatedAt = member.ActivatedAt.String()
	}

	return &domain.OrganizationMember{
		PkId:             member.Pkid,
		OrganizationPkID: member.OrganizationPkid,
		UserPkID:         member.UserPkid,
		Role:             member.Role,
		User:             mapUserModelToDomain(member.User),
		ActivatedAt:      activatedAt,
		CreatedAt:        member.CreatedAt.String(),
		UpdatedAt:        member.UpdatedAt.String(),
	}, nil
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

func (r *OrganizationRepository) GetOwnerOrgByPkId(ctx context.Context, ownerID, pkId int64) (*domain.Organization, *domain.Error) {
	var org model.Organization

	err := r.store.DB().Where("owner_id = ? AND pkid = ?", ownerID, pkId).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrOrgNotFound
		}

		return nil, domain.ErrDatabaseQuery
	}

	return &domain.Organization{
		PkId:    org.Pkid,
		OwnerID: org.OwnerID,
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

	// -- Transaction
	tx, doneTx := r.store.NewTransaction()
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

	return mapOrg(newOrg, []domain.OrganizationMember{
		*mapOrgMemberUserModelToDomain(ownerMember, owner),
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

func (r *OrganizationRepository) AddMemberToOrg(ctx context.Context, orgPkID int64, userPkID *int64, role string) (*domain.OrganizationMember, *domain.Error) {
	var newMember = model.OrganizationMember{
		OrganizationPkid: orgPkID,
		UserPkid:         userPkID,
		Role:             role,
	}
	err := r.store.DB().Create(&newMember).Error
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	var user *domain.User

	if newMember.UserPkid != nil {
		user, _ = r.userRepository.GetUserByPkID(context.Background(), *newMember.UserPkid)
	}

	return mapOrgMemberUserModelToDomain(newMember, user), nil
}

func (r *OrganizationRepository) SetOrgMemberActivatedAt(ctx context.Context, pkID int64, activatedAt time.Time) (*domain.OrganizationMember, *domain.Error) {
	var member model.OrganizationMember

	err := r.store.DB().Model(&member).Clauses(clause.Returning{}).Where("pkid = ?", pkID).Update("activated_at", activatedAt).Error
	if err != nil {
		return nil, domain.ErrDatabaseMutation
	}

	return mapOrgMemberUserModelToDomain(member, nil), nil
}

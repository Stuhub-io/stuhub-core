package organizationutils

import (
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
	"github.com/Stuhub-io/utils/userutils"
)

type MemberWithUser struct {
	model.OrganizationMember
	User model.User `gorm:"foreignKey:user_pkid" json:"user"` // Define foreign key relationship
}

type OrganizationWithMembers struct {
	model.Organization
	Owner   model.User       `gorm:"foreignKey:owner_id"          json:"owner"`
	Members []MemberWithUser `gorm:"foreignKey:organization_pkid" json:"members"` // Consider JSON tag for future use
}

func TransformOrganizationMemberModelToDomain(member MemberWithUser) *domain.OrganizationMember {
	activatedAt := ""
	if member.ActivatedAt != nil {
		activatedAt = member.ActivatedAt.String()
	}

	return &domain.OrganizationMember{
		PkID:             member.Pkid,
		OrganizationPkID: member.OrganizationPkid,
		UserPkID:         member.UserPkid,
		Role:             member.Role,
		User:             userutils.TransformUserModelToDomain(member.User),
		ActivatedAt:      activatedAt,
		CreatedAt:        member.CreatedAt.String(),
		UpdatedAt:        member.UpdatedAt.String(),
	}
}

func TransformOrganizationMemberModelToDomain_Many(members []MemberWithUser) []domain.OrganizationMember {
	domainMembers := make([]domain.OrganizationMember, 0, len(members))
	for _, member := range members {
		domainMembers = append(domainMembers, *TransformOrganizationMemberModelToDomain(member))
	}

	return domainMembers
}

func TransformOrganizationModelToDomain(model OrganizationWithMembers) *domain.Organization {
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
		Owner:       userutils.TransformUserModelToDomain(model.Owner),
		Members:     TransformOrganizationMemberModelToDomain_Many(model.Members),
	}
}

func TransformOrganizationModelToDomain_Many(models []OrganizationWithMembers) []*domain.Organization {
	domainOrgs := make([]*domain.Organization, 0, len(models))
	for _, org := range models {
		domainOrg := TransformOrganizationModelToDomain(org)
		domainOrgs = append(domainOrgs, domainOrg)
	}

	return domainOrgs
}

func TransformOrganizationModelToDomain_New(org model.Organization, ownerMember model.OrganizationMember, owner *domain.User) *domain.Organization {
	activatedAt := ""
	if ownerMember.ActivatedAt != nil {
		activatedAt = ownerMember.ActivatedAt.String()
	}

	member := domain.OrganizationMember{
		PkID:             ownerMember.Pkid,
		OrganizationPkID: ownerMember.OrganizationPkid,
		UserPkID:         ownerMember.UserPkid,
		Role:             ownerMember.Role,
		User:             owner,
		ActivatedAt:      activatedAt,
		CreatedAt:        ownerMember.CreatedAt.String(),
		UpdatedAt:        ownerMember.UpdatedAt.String(),
	}

	return &domain.Organization{
		ID:          org.ID,
		PkId:        org.Pkid,
		OwnerID:     org.OwnerID,
		Name:        org.Name,
		Slug:        org.Slug,
		Description: org.Description,
		Avatar:      org.Avatar,
		CreatedAt:   org.CreatedAt.String(),
		UpdatedAt:   org.UpdatedAt.String(),
		Members:     []domain.OrganizationMember{member},
	}
}

func TransformOrganizationMemberModelToDomain_New(member model.OrganizationMember, user *domain.User) *domain.OrganizationMember {
	activatedAt := ""
	if member.ActivatedAt != nil {
		activatedAt = member.ActivatedAt.String()
	}

	return &domain.OrganizationMember{
		PkID:             member.Pkid,
		OrganizationPkID: member.OrganizationPkid,
		UserPkID:         member.UserPkid,
		Role:             member.Role,
		User:             user,
		ActivatedAt:      activatedAt,
		CreatedAt:        member.CreatedAt.String(),
		UpdatedAt:        member.UpdatedAt.String(),
	}
}

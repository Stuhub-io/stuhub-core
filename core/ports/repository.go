package ports

import (
	"context"

	"github.com/Stuhub-io/core/domain"
)

type UserRepository interface {
	GetByID(ctx context.Context, id string) (*domain.User, *domain.Error)
	GetUserByPkID(ctx context.Context, pkID int64) (*domain.User, *domain.Error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, *domain.Error)
	GetOrCreateUserByEmail(ctx context.Context, email, salt string) (*domain.User, *domain.Error)
	CreateUserWithGoogleInfo(ctx context.Context, email, salt, firstName, lastName, avatar string) (*domain.User, *domain.Error)
	SetUserPassword(ctx context.Context, PkID int64, hashedPassword string) *domain.Error
	CheckPassword(ctx context.Context, email, rawPassword string, hasher Hasher) (bool, *domain.Error)
	UpdateUserInfo(ctx context.Context, PkID int64, firstName, lastName, avatar string) (*domain.User, *domain.Error)
}

type OrganizationRepository interface {
	GetOrgMembers(ctx context.Context, pkID int64) ([]*domain.OrganizationMember, *domain.Error)
	GetOrgBySlug(ctx context.Context, slug string) (*domain.Organization, *domain.Error)
	GetOwnerOrgByName(ctx context.Context, ownerID int64, name string) (*domain.Organization, *domain.Error)
	GetOwnerOrgByPkId(ctx context.Context, ownerID, pkId int64) (*domain.Organization, *domain.Error)
	GetOrgsByUserPkID(ctx context.Context, usePkID int64) ([]*domain.Organization, *domain.Error)
	GetOrgMemberByEmail(ctx context.Context, orgPkId int64, email string) (*domain.OrganizationMember, *domain.Error)
	CreateOrg(ctx context.Context, userPkID int64, name, description, avatar string) (*domain.Organization, *domain.Error)
	AddMemberToOrg(ctx context.Context, orgPkID int64, userPkID *int64, role string) (*domain.OrganizationMember, *domain.Error)
}

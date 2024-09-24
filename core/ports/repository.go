package ports

import (
	"context"
	"time"

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
	SetUserActivatedAt(ctx context.Context, pkID int64, activatedAt time.Time) (*domain.User, *domain.Error)
}

type OrganizationRepository interface {
	GetOrgMembers(ctx context.Context, pkID int64) ([]*domain.OrganizationMember, *domain.Error)
	GetOrgBySlug(ctx context.Context, slug string) (*domain.Organization, *domain.Error)
	GetOwnerOrgByName(ctx context.Context, ownerPkID int64, name string) (*domain.Organization, *domain.Error)
	GetOwnerOrgByPkId(ctx context.Context, ownerPkID, pkId int64) (*domain.Organization, *domain.Error)
	GetOrgsByUserPkID(ctx context.Context, usePkID int64) ([]*domain.Organization, *domain.Error)
	GetOrgMemberByEmail(ctx context.Context, orgPkId int64, email string) (*domain.OrganizationMember, *domain.Error)
	GetOrgMemberByUserPkID(ctx context.Context, orgPkId int64, userPkId int64) (*domain.OrganizationMember, *domain.Error)
	CreateOrg(ctx context.Context, userPkID int64, name, description, avatar string) (*domain.Organization, *domain.Error)
	AddMemberToOrg(ctx context.Context, orgPkID int64, userPkID *int64, role string) (*domain.OrganizationMember, *domain.Error)
	SetOrgMemberActivatedAt(ctx context.Context, pkID int64, activatedAt time.Time) (*domain.OrganizationMember, *domain.Error)
}

type SpaceRepository interface {
	CreateSpace(ctx context.Context, orgPkID int64, ownerPkID int64, isPrivate bool, name, description string) (*domain.Space, *domain.Error)
	GetSpacesByOrgPkID(ctx context.Context, orgPkID int64) ([]domain.Space, *domain.Error)
}

type PageRepository interface {
	CreatePage(ctx context.Context, spacePkID int64, name string, viewType domain.PageViewType, ParentPagePkID *int64) (*domain.Page, *domain.Error)
	GetPagesBySpacePkID(ctx context.Context, spacePkID int64) ([]domain.Page, *domain.Error)
	DeletePageByPkID(ctx context.Context, pagePkID int64, userPkID int64) (*domain.Page, *domain.Error)
	GetPageByID(ctx context.Context, pageID string) (*domain.Page, *domain.Error)
	UpdatePageByID(ctx context.Context, pageID string, page domain.PageInput) (*domain.Page, *domain.Error)
}

type DocumentRepository interface {
	CreateDocument(ctx context.Context, pagePkID int64, JsonContent string) (*domain.Document, *domain.Error)
	GetDocumentByPagePkID(ctx context.Context, pagePkID int64) (*domain.Document, *domain.Error)
	GetDocumentByPkID(ctx context.Context, pkID int64) (*domain.Document, *domain.Error)
	UpdateDocument(ctx context.Context, pagePkID int64, content string) (*domain.Document, *domain.Error)
}

type OrganizationInviteRepository interface {
	CreateInvite(ctx context.Context, organizationPkID int64, userPkID int64) (*domain.OrganizationInvite, *domain.Error)
	GetInviteByID(ctx context.Context, inviteID string) (*domain.OrganizationInvite, *domain.Error)
}

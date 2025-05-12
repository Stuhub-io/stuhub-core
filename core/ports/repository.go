package ports

import (
	"context"
	"time"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/internal/repository/model"
)

type UserRepository interface {
	Search(
		ctx context.Context,
		query domain.UserSearchQuery,
		currentUser *domain.User,
	) ([]domain.User, *domain.Error)
	GetByID(ctx context.Context, id string) (*domain.User, *domain.Error)
	GetUserByPkID(ctx context.Context, pkID int64) (*domain.User, *domain.Error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, *domain.Error)
	GetOrCreateUserByEmail(ctx context.Context, email, salt string) (*domain.User, *domain.Error, bool) // bool - iscreated
	CreateUserWithGoogleInfo(
		ctx context.Context,
		email, salt, firstName, lastName, avatar string,
	) (*domain.User, *domain.Error)
	SetUserPassword(ctx context.Context, PkID int64, hashedPassword string) *domain.Error
	CheckPassword(
		ctx context.Context,
		email, rawPassword string,
		hasher Hasher,
	) (bool, *domain.Error)
	UpdateUserInfo(
		ctx context.Context,
		PkID int64,
		firstName, lastName, avatar string,
	) (*domain.User, *domain.Error)
	SetUserActivatedAt(
		ctx context.Context,
		pkID int64,
		activatedAt time.Time,
	) (*domain.User, *domain.Error)
	UnsafeListUsers(
		ctx context.Context,
		query domain.UserListQuery,
	) ([]domain.User, *domain.Error)
}

type OrganizationRepository interface {
	GetOrgMembers(ctx context.Context, pkID int64) ([]domain.OrganizationMember, *domain.Error)
	GetOrgByPkID(ctx context.Context, pkID int64) (*domain.Organization, *domain.Error)
	GetOrgBySlug(ctx context.Context, slug string) (*domain.Organization, *domain.Error)
	GetOwnerOrgByName(
		ctx context.Context,
		ownerPkID int64,
		name string,
	) (*domain.Organization, *domain.Error)
	GetOwnerOrgByPkID(
		ctx context.Context,
		ownerPkID, pkId int64,
	) (*domain.Organization, *domain.Error)
	GetOrgsByUserPkID(ctx context.Context, usePkID int64) ([]*domain.Organization, *domain.Error)
	GetOrgMemberByEmail(
		ctx context.Context,
		orgPkID int64,
		email string,
	) (*domain.OrganizationMember, *domain.Error)
	GetOrgMemberByUserPkID(
		ctx context.Context,
		orgPkID int64,
		userPkID int64,
	) (*domain.OrganizationMember, *domain.Error)
	CreateOrg(
		ctx context.Context,
		userPkID int64,
		name, description, avatar string,
	) (*domain.Organization, *domain.Error)
	AddMemberToOrg(
		ctx context.Context,
		orgPkID int64,
		userPkID *int64,
		role string,
	) (*domain.OrganizationMember, *domain.Error)
	SetOrgMemberActivatedAt(
		ctx context.Context,
		pkID int64,
		activatedAt time.Time,
	) (*domain.OrganizationMember, *domain.Error)
}

type PageRepository interface {
	List(
		ctx context.Context,
		query domain.PageListQuery,
		curUser *domain.User,
	) ([]domain.Page, *domain.Error)
	Update(
		ctx context.Context,
		pagePkID int64,
		page domain.PageUpdateInput,
	) (*domain.Page, *domain.Error)
	Move(ctx context.Context, pagePkID int64, parentPagePkID *int64) (*domain.Page, *domain.Error)
	GetByID(
		ctx context.Context,
		pageID string,
		pagePkID *int64,
		detailOption domain.PageDetailOptions,
		actorPkID *int64,
	) (*domain.Page, *domain.Error)
	Archive(ctx context.Context, pagePkID int64) (*domain.Page, *domain.Error)
	UpdateGeneralAccess(
		ctx context.Context,
		pagePkID int64,
		updateInput domain.PageGeneralAccessUpdateInput,
	) (*domain.Page, *domain.Error)
	// publication
	CreatePublicToken(ctx context.Context, pagePkID int64) (*domain.PagePublicToken, *domain.Error)
	ArchiveAllPublicToken(ctx context.Context, pagePkID int64) *domain.Error
	GetPublicTokenByID(
		ctx context.Context,
		publicTokenID string,
	) (*domain.PagePublicToken, *domain.Error)

	// Document Page
	CreateDocumentPage(
		ctx context.Context,
		page domain.DocumentPageInput,
	) (*domain.Page, *domain.Error)
	UpdateContent(
		ctx context.Context,
		pagePkID int64,
		content domain.DocumentInput,
	) (*domain.Page, *domain.Error)

	// Asset Page
	CreateAsset(ctx context.Context, asset domain.AssetPageInput) (*domain.Page, *domain.Error)

	// Page Role
	CreatePageRole(
		ctx context.Context,
		createInput domain.PageRoleCreateInput,
	) (*domain.PageRoleUser, *domain.Error)
	GetPageRoleByEmail(
		ctx context.Context,
		pagePkID int64, email string,
	) (*domain.PageRoleUser, *domain.Error)
	GetPageRoles(
		ctx context.Context,
		pagePkID int64,
	) ([]domain.PageRoleUser, *domain.Error)
	GetPagesRole(
		ctx context.Context,
		input domain.PageRolePermissionBatchCheckInput,
	) (permissions []domain.PageRolePermissionCheckInput, err *domain.Error)
	UpdatePageRole(
		ctx context.Context,
		updateInput domain.PageRoleUpdateInput,
	) *domain.Error
	DeletePageRole(
		ctx context.Context,
		updateInput domain.PageRoleDeleteInput,
	) *domain.Error
	CheckPermission(
		ctx context.Context,
		input domain.PageRolePermissionCheckInput,
	) domain.PageRolePermissions

	SyncPageRoleWithNewUser(
		ctx context.Context,
		user domain.User,
	) *domain.Error

	// Page Role Request
	CreatePageAccessRequest(
		ctx context.Context,
		createInput domain.PageRoleRequestCreateInput,
	) (*domain.PageRoleRequestLog, *domain.Error)
	ListPageAccessRequestByPagePkID(
		ctx context.Context,
		q domain.PageRoleRequestLogQuery,
	) ([]domain.PageRoleRequestLog, *domain.Error)
	UpdatePageAccessRequestStatus(
		ctx context.Context,
		q domain.PageRoleRequestLogQuery,
		status domain.PageRoleRequestLogStatus,
	) *domain.Error

	// Page Star
	StarPage(ctx context.Context, input domain.StarPageInput) (*domain.PageStar, *domain.Error)
	UnstarPage(ctx context.Context, input domain.StarPageInput) *domain.Error
}

type OrganizationInviteRepository interface {
	CreateInvite(
		ctx context.Context,
		organizationPkId int64,
		userPkId int64,
	) (*domain.OrganizationInvite, *domain.Error)
	UpdateInvite(
		ctx context.Context,
		invite model.OrganizationInvite,
	) (*domain.OrganizationInvite, *domain.Error)
	GetInviteByID(ctx context.Context, inviteID string) (*domain.OrganizationInvite, *domain.Error)
}

type PageAccessLogRepository interface {
	GetByUserPKID(
		ctx context.Context,
		query domain.CursorPagination[time.Time],
		userPkID int64,
	) ([]domain.PageAccessLog, *domain.Error)
	Upsert(
		ctx context.Context,
		pagePkID,
		userPkID int64,
		action domain.PageAccessAction,
	) (int64, *domain.Error)
}

type ActivityRepository interface {
	List(
		ctx context.Context,
		query domain.ActivityListQuery,
	) ([]domain.Activity, *domain.Error)
	Create(
		ctx context.Context,
		input domain.ActivityInput,
	) (*domain.Activity, *domain.Error)
}

type ActivityV2Repository interface {
	List(ctx context.Context, query domain.ActivityV2ListQuery) ([]domain.ActivityV2, *domain.Error)
	Create(ctx context.Context, input domain.ActivityV2Input) (*domain.ActivityV2, *domain.Error)
	Update(ctx context.Context, activityPkID int64, input domain.ActivityV2Input) (*domain.ActivityV2, *domain.Error)
	One(ctx context.Context, query domain.ActivityV2ListQuery) (*domain.ActivityV2, *domain.Error)
}

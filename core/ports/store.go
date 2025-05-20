package ports

import (
	"github.com/Stuhub-io/core/domain"
)

type DBStore interface {
	Cache() CacheStore
	NewTransaction() (*Repository, IFinallyFunc)
	Shutdown() *domain.Error
}

type IFinallyFunc interface {
	Commit() *domain.Error
	Rollback(*domain.Error) *domain.Error
}

type Repository struct {
	Store              DBStore
	User               UserRepository
	Organization       OrganizationRepository
	Page               PageRepository
	OrganizationInvite OrganizationInviteRepository
	PageAccessLog      PageAccessLogRepository
	Activity           ActivityRepository
	ActivityV2         ActivityV2Repository
}

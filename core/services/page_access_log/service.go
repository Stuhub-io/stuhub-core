package pageAccessLog

import (
	"context"
	"slices"
	"sort"
	"time"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/utils/pageutils"
	sliceutils "github.com/Stuhub-io/utils/slice"
)

type Service struct {
	repo *ports.Repository
}

type NewServiceParams struct {
	*ports.Repository
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		repo: params.Repository,
	}
}

func (s *Service) GetLogsByUser(
	query domain.CursorPagination[time.Time],
	user *domain.User,
) ([]domain.PageAccessLog, *time.Time, *domain.Error) {
	logs, err := s.repo.PageAccessLog.GetByUserPKID(context.Background(), query, user.PkID)
	if err != nil {
		return nil, nil, err
	}

	nextCursor := domain.CalculateNextCursor[domain.PageAccessLog, time.Time](query.Limit, logs, "LastAccessed")

	flatPages := sliceutils.FlatMap(logs, func(log domain.PageAccessLog) []domain.Page {
		return append(log.ParentPages, log.Page)
	})
	flatPages = sliceutils.UniqueByField(flatPages, "PkID")

	permissionInputs, _ := s.repo.Page.GetPagesRole(
		context.Background(),
		domain.PageRolePermissionBatchCheckInput{
			User:  user,
			Pages: flatPages,
		},
	)

	permissionsMapper := map[int64]domain.PageRolePermissions{}
	inheritRolePagesMapper := map[int64]string{}

	// get permissions for not inherit pages
	for _, page := range flatPages {
		var pageRole *domain.PageRole

		foundPageInPermission := sliceutils.Find(
			permissionInputs,
			func(p domain.PageRolePermissionCheckInput) bool {
				return p.Page.PkID == page.PkID
			},
		)
		if foundPageInPermission != nil {
			pageRole = foundPageInPermission.PageRole
		}

		if pageRole != nil && pageRole.String() == domain.PageInherit.String() && page.Path != "" {
			inheritRolePagesMapper[page.PkID] = page.Path
			continue
		}

		permissions := s.repo.Page.CheckPermission(
			context.Background(),
			domain.PageRolePermissionCheckInput{
				User:     user,
				Page:     page,
				PageRole: pageRole,
			},
		)

		permissionsMapper[page.PkID] = permissions
	}

	// get permissions for inherit pages
	findInheritPermissions := func(parentPkIds []int64) domain.PageRolePermissions {
		slices.Reverse(parentPkIds)
		for _, pkID := range parentPkIds {
			if permissions, ok := permissionsMapper[pkID]; ok {
				if permissions.CanView {
					return permissions
				}
			}
		}
		return domain.PageRolePermissions{
			CanView: false,
		}
	}

	for pkID, path := range inheritRolePagesMapper {
		permissionsMapper[pkID] = findInheritPermissions(pageutils.PagePathToPkIDs(path))
	}

	// filter logs after checking permission
	for i := 0; i < len(logs); i++ {
		log := &logs[i]

		permissions, ok := permissionsMapper[log.Page.PkID]
		if !ok || !permissions.CanView {
			logs = sliceutils.Filter(logs, func(l domain.PageAccessLog) bool {
				return l.PkID != log.PkID
			})
			continue
		}

		log.Page.Permissions = &permissions
		log.IsShared = !log.Page.IsAuthor(user.PkID)

		for _, page := range log.ParentPages {
			permissions, ok := permissionsMapper[page.PkID]
			if !ok || !permissions.CanView {
				log.ParentPages = sliceutils.Filter(log.ParentPages, func(p domain.Page) bool {
					return p.PkID != page.PkID
				})
			}
		}

		sort.Slice(log.ParentPages, func(i, j int) bool {
			return log.ParentPages[i].PkID < log.ParentPages[j].PkID
		})
	}

	return logs, nextCursor, err
}

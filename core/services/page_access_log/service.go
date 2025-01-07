package pageAccessLog

import (
	"context"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	sliceutils "github.com/Stuhub-io/utils/slice"
)

type Service struct {
	pageRepository          ports.PageRepository
	pageAccessLogRepository ports.PageAccessLogRepository
}

type NewServiceParams struct {
	ports.PageRepository
	ports.PageAccessLogRepository
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		pageRepository:          params.PageRepository,
		pageAccessLogRepository: params.PageAccessLogRepository,
	}
}

func (s *Service) GetLogsByUser(
	query domain.OffsetBasedPagination,
	user *domain.User,
) ([]domain.PageAccessLog, *domain.Error) {
	logs, err := s.pageAccessLogRepository.GetByUserPKID(context.Background(), query, user.PkID)

	flatPages := sliceutils.FlatMap(logs, func(log domain.PageAccessLog) []domain.Page {
		return append(log.ParentPages, log.Page)
	})
	flatPages = sliceutils.UniqueByField(flatPages, "PkID")

	permissionInputs, _ := s.pageRepository.GetPagesRole(
		context.Background(),
		domain.PageRolePermissionBatchCheckInput{
			User:  user,
			Pages: flatPages,
		},
	)

	permissionsMapper := map[int64]domain.PageRolePermissions{}

	// checking permission for all pages
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

		permissions := s.pageRepository.CheckPermission(
			context.Background(),
			domain.PageRolePermissionCheckInput{
				User:     user,
				Page:     page,
				PageRole: pageRole,
			},
		)
		if !permissions.CanView {
			flatPages = sliceutils.Filter(flatPages, func(p domain.Page) bool {
				return p.PkID != page.PkID
			})
			continue
		}

		permissionsMapper[page.PkID] = permissions
	}

	// filter logs after checking permission
	for i := 0; i < len(logs); i++ {
		log := &logs[i]

		validPage := sliceutils.Find(flatPages, func(p domain.Page) bool {
			return p.PkID == log.Page.PkID
		})
		if validPage == nil {
			logs = sliceutils.Filter(logs, func(l domain.PageAccessLog) bool {
				return l.PkID != log.PkID
			})
			continue
		}

		permission := permissionsMapper[log.Page.PkID]
		log.Page.Permissions = &permission

		for _, page := range log.ParentPages {
			validParentPage := sliceutils.Find(flatPages, func(p domain.Page) bool {
				return p.PkID == page.PkID
			})
			if validParentPage == nil {
				log.ParentPages = sliceutils.Filter(log.ParentPages, func(p domain.Page) bool {
					return p.PkID != page.PkID
				})
			}
		}
	}

	return logs, err
}

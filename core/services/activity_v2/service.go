package activity_v2

import (
	"context"
	"fmt"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/logger"
	sliceutils "github.com/Stuhub-io/utils/slice"
)

type Service struct {
	cfg                  config.Config
	logger               logger.Logger
	pageRepository       ports.PageRepository
	activityV2Repository ports.ActivityV2Repository
	userRepository       ports.UserRepository
	orgRepository        ports.OrganizationRepository
}

type NewServiceParams struct {
	config.Config
	logger.Logger
	ports.PageRepository
	ports.ActivityV2Repository
	ports.OrganizationRepository
	ports.UserRepository
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		cfg:                  params.Config,
		logger:               params.Logger,
		pageRepository:       params.PageRepository,
		activityV2Repository: params.ActivityV2Repository,
		orgRepository:        params.OrganizationRepository,
		userRepository:       params.UserRepository,
	}
}

func (s Service) ListPageActivities(curUser *domain.User, pagePkID int64, pagination domain.CursorPagination[time.Time]) ([]domain.ActivityV2, *domain.Error) {

	// CHECK PERMISSION
	if curUser == nil {
		return nil, domain.ErrUnauthorized
	}

	page, err := s.pageRepository.GetByID(context.Background(), "", &pagePkID, domain.PageDetailOptions{}, &curUser.PkID)
	if err != nil {
		e := fmt.Errorf(err.Message)
		s.logger.Error(e, "[Activity]: "+err.Message)
		return nil, domain.ErrNotFound
	}

	var userRole *domain.PageRole = nil
	pageRole, _ := s.pageRepository.GetPageRoleByEmail(context.Background(), pagePkID, curUser.Email)
	if pageRole != nil {
		userRole = &pageRole.Role
	}

	permisisons := s.pageRepository.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *page,
		User:     curUser,
		PageRole: userRole,
	})

	if !permisisons.CanView {
		return nil, domain.ErrPermissionDenied
	}

	// QUERY ACTIVITIES FROM LOG DB
	// List all viewable pages
	pages, err := s.pageRepository.List(
		context.Background(),
		domain.PageListQuery{
			IsAll:          true,
			ParentPagePkID: &page.PkID,
		},
		curUser,
	)
	if err != nil {
		return nil, domain.ErrBadRequest
	}

	AllChildrenPagePkIDs := sliceutils.Map(pages, func(page domain.Page) int64 {
		return page.PkID
	})

	pagePkIds := append([]int64{page.PkID}, AllChildrenPagePkIDs...)

	activities, err := s.activityV2Repository.List(context.Background(), domain.ActivityV2ListQuery{
		RelatedPagePkIDs: pagePkIds,
		Limit:            &pagination.Limit,
		EndTime:          &pagination.Cursor,
	})

	if err != nil {
		e := fmt.Errorf(err.Message)
		s.logger.Errorf(e, "[Activity]: Get Activities From PagePkIDs Error")
		return nil, domain.ErrDatabaseQuery
	}

	return activities, nil
}

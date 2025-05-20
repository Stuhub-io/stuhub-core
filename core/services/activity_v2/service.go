package activity_v2

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/logger"
	"github.com/Stuhub-io/utils/pageutils"
	sliceutils "github.com/Stuhub-io/utils/slice"
)

type Service struct {
	cfg    config.Config
	logger logger.Logger
	repo   *ports.Repository
}

type NewServiceParams struct {
	config.Config
	logger.Logger
	*ports.Repository
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		cfg:    params.Config,
		logger: params.Logger,
		repo:   params.Repository,
	}
}

func (s Service) ListPageActivities(curUser *domain.User, pagePkID int64, pagination domain.CursorPagination[time.Time]) ([]domain.ActivityV2, *domain.Error) {

	// CHECK PERMISSION
	if curUser == nil {
		return nil, domain.ErrUnauthorized
	}

	page, err := s.repo.Page.GetByID(context.Background(), "", &pagePkID, domain.PageDetailOptions{}, &curUser.PkID)
	if err != nil {
		e := fmt.Errorf(err.Message)
		s.logger.Error(e, "[Activity]: "+err.Message)
		return nil, domain.ErrNotFound
	}

	var userRole *domain.PageRole = nil
	pageRole, _ := s.repo.Page.GetPageRoleByEmail(context.Background(), pagePkID, curUser.Email)
	if pageRole != nil {
		userRole = &pageRole.Role
	}

	permisisons := s.repo.Page.CheckPermission(context.Background(), domain.PageRolePermissionCheckInput{
		Page:     *page,
		User:     curUser,
		PageRole: userRole,
	})

	if !permisisons.CanView {
		return nil, domain.ErrPermissionDenied
	}

	// QUERY ACTIVITIES FROM LOG DB
	// List all viewable pages
	pages, err := s.repo.Page.List(
		context.Background(),
		domain.PageListQuery{
			IsAll:         true,
			PathBeginWith: pageutils.AppendPath(page.Path, strconv.FormatInt(pagePkID, 10)),
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

	activities, err := s.repo.ActivityV2.List(context.Background(), domain.ActivityV2ListQuery{
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

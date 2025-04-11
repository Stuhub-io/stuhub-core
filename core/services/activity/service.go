package activity

import (
	"context"
	"fmt"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/logger"
	sliceutils "github.com/Stuhub-io/utils/slice"
)

type Service struct {
	cfg                config.Config
	logger             logger.Logger
	pageRepository     ports.PageRepository
	activityRepository ports.ActivityRepository
	userRepository     ports.UserRepository
	orgRepository      ports.OrganizationRepository
}

type NewServiceParams struct {
	config.Config
	logger.Logger
	ports.PageRepository
	ports.ActivityRepository
	ports.OrganizationRepository
	ports.UserRepository
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		cfg:                params.Config,
		logger:             params.Logger,
		pageRepository:     params.PageRepository,
		activityRepository: params.ActivityRepository,
		orgRepository:      params.OrganizationRepository,
		userRepository:     params.UserRepository,
	}
}

func (s Service) TrackUserVisitPage(curUser *domain.User, pagePkID int64) *domain.Error {
	if curUser == nil {
		return domain.ErrUnauthorized
	}

	p, e := s.pageRepository.GetByID(context.Background(), "", &pagePkID, domain.PageDetailOptions{}, nil)

	if e != nil {
		return e
	}

	label := "User Visited Page"
	_, er := s.activityRepository.Create(context.Background(), domain.ActivityInput{
		ActionCode: domain.ActionUserVisitPage,
		ActorPkID:  curUser.PkID,
		PagePkID:   &p.PkID,
		OrgPkID:    &p.OrganizationPkID,
		Label:      &label,
	})
	if er != nil {
		return er
	}
	return nil
}

func (s Service) TrackUserVisitOrganization(curUser *domain.User, orgPkID int64) *domain.Error {
	if curUser == nil {
		return domain.ErrUnauthorized
	}
	org, e := s.orgRepository.GetOrgByPkID(context.Background(), orgPkID)

	if e != nil {
		return e
	}

	label := "User Visited Page"

	_, er := s.activityRepository.Create(context.Background(), domain.ActivityInput{
		ActionCode: domain.ActionUserVisitPage,
		ActorPkID:  curUser.PkID,
		PagePkID:   nil,
		OrgPkID:    &org.PkId,
		Label:      &label,
	})
	if er != nil {
		return er
	}

	return nil
}

func (s Service) ListPageActivities(curUser *domain.User, pagePkID int64) ([]domain.Activity, *domain.Error) {

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

	childPagePkIds := sliceutils.Map(pages, func(page domain.Page) int64 {
		return page.PkID
	})

	pagePkIds := append([]int64{page.PkID}, childPagePkIds...)

	activities, err := s.activityRepository.List(context.Background(), domain.ActivityListQuery{
		PagePkIDs: pagePkIds,
	})

	if err != nil {
		e := fmt.Errorf(err.Message)
		s.logger.Errorf(e, "[Activity]: Get Activities From PagePkIDs Error")
		return nil, domain.ErrDatabaseQuery
	}

	detail_activities, err := s.EnrichActivities(activities, curUser)

	if err != nil {
		return nil, domain.ErrBadRequest
	}

	return detail_activities, nil
}

func (s Service) EnrichActivities(activities []domain.Activity, curUser *domain.User) ([]domain.Activity, *domain.Error) {
	pagePkIDs := make([]int64, 0, len(activities))
	actorPkIDs := make([]int64, 0, len(activities))

	pagePkIDsMap := make(map[int64]bool)
	actorPkIDsMap := make(map[int64]bool)

	for _, activity := range activities {
		if activity.PagePkID != nil {
			if !pagePkIDsMap[*activity.PagePkID] {
				pagePkIDsMap[*activity.PagePkID] = true
				pagePkIDs = append(pagePkIDs, *activity.PagePkID)
			}
		}

		if !actorPkIDsMap[activity.ActorPkID] {
			actorPkIDsMap[activity.ActorPkID] = true
			actorPkIDs = append(actorPkIDs, activity.ActorPkID)
		}
	}

	pages, err := s.pageRepository.List(context.Background(), domain.PageListQuery{
		PagePkIDs: pagePkIDs,
		IsAll:     true,
	}, curUser)

	if err != nil {
		return nil, err
	}

	pagesMap := make(map[int64]domain.Page)
	for _, page := range pages {
		pagesMap[page.PkID] = page
	}

	users, err := s.userRepository.UnsafeListUsers(context.Background(), domain.UserListQuery{
		UserPkIDs: actorPkIDs,
	})

	if err != nil {
		e := fmt.Errorf(err.Message)
		s.logger.Error(e, "[Activity]: Unsafe List User error")
		return nil, domain.ErrDatabaseQuery
	}

	usersMap := make(map[int64]domain.User)
	for _, user := range users {
		usersMap[user.PkID] = user
	}

	for i := range activities {
		actor := usersMap[activities[i].ActorPkID]
		activities[i].Actor = &actor
		if activities[i].PagePkID != nil {
			page := pagesMap[*activities[i].PagePkID]
			activities[i].Page = &page
		}
	}

	return activities, nil
}

package api

import (
	"net/http"
	"time"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/services/activity"
	"github.com/Stuhub-io/internal/api/decorators"
	"github.com/Stuhub-io/internal/api/middleware"
	"github.com/Stuhub-io/internal/api/request"
	"github.com/Stuhub-io/internal/api/response"
	"github.com/Stuhub-io/utils/organizationutils"
	"github.com/Stuhub-io/utils/pageutils"
	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	activityService *activity.Service
	authMiddleware  *middleware.AuthMiddleware
}

type NewActivityHandlerParams struct {
	Router          *gin.RouterGroup
	ActivityService *activity.Service
	AuthMiddleware  *middleware.AuthMiddleware
}

func UseActivityHandler(params NewActivityHandlerParams) {
	handler := &ActivityHandler{
		activityService: params.ActivityService,
		authMiddleware:  params.AuthMiddleware,
	}

	router := params.Router.Group("/activity-services")
	authMiddleware := params.AuthMiddleware

	router.Use(authMiddleware.Authenticated())
	router.POST("/pages/:"+pageutils.PagePkIDParam+"/track-visit", decorators.RequiredAuth(decorators.CurrentUser(handler.TrackUserVisitPage)))
	router.POST("/orgs/:"+organizationutils.OrgPkIDParam+"/track-visit", decorators.RequiredAuth(decorators.CurrentUser(handler.TrackUserVisitOrg)))
	router.GET("/pages/:"+pageutils.PagePkIDParam+"/activities", decorators.RequiredAuth(decorators.CurrentUser(handler.ListActivities)))
}

func (h *ActivityHandler) TrackUserVisitPage(c *gin.Context, curUser *domain.User) {
	pagePkID, ok := pageutils.GetPagePkIDParam(c)
	if !ok {
		response.BindError(c, "pagePkID is missing or invalid")
		return
	}
	e := h.activityService.TrackUserVisitPage(curUser, pagePkID)
	if e != nil {
		response.WithErrorMessage(c, e.Code, e.Error, e.Message)
		return
	}
	response.WithMessage(c, 200, "User visit page tracked successfully")
}

func (h *ActivityHandler) TrackUserVisitOrg(c *gin.Context, curUser *domain.User) {
	orgPkID, ok := organizationutils.GetOrgPkIDParam(c)
	if !ok {
		response.BindError(c, "orgPkID is missing or invalid")
		return
	}

	e := h.activityService.TrackUserVisitOrganization(curUser, orgPkID)
	if e != nil {
		response.WithErrorMessage(c, e.Code, e.Error, e.Message)
		return
	}
	response.WithMessage(c, 200, "User visit org tracked successfully")
}

func (h *ActivityHandler) ListActivities(c *gin.Context, curUser *domain.User) {
	pagePkID, ok := pageutils.GetPagePkIDParam(c)
	if !ok {
		response.BindError(c, "pagePkID is missing or invalid")
		return
	}

	var query request.ActivityPaginationRequest
	if verr := request.Validate(c, &query); verr != nil {
		response.BindError(c, verr.Error())
		return
	}

	if verr := request.Validate(c, &query); verr != nil {
		response.BindError(c, verr.Error())
		return
	}

	cursor := time.Now()
	if query.EndTime != "" {
		time, err := time.Parse(time.RFC3339, query.EndTime)
		if err != nil {
			response.BindError(c, err.Error())
			return
		}
		cursor = time
	}

	activities, err := h.activityService.ListPageActivities(curUser, pagePkID, domain.CursorPagination[time.Time]{
		Limit:  query.Limit,
		Cursor: cursor,
	})

	nextCursorStr := domain.CalculateNextCursor[domain.Activity, string](1, activities, "CreatedAt")
	nextCursor := time.Unix(0, 0)

	if nextCursorStr != nil {
		time, err := time.Parse(time.RFC3339, *nextCursorStr)
		if err != nil {
			response.BindError(c, err.Error())
			return
		}
		nextCursor = time
	}

	if err != nil {
		response.BindError(c, err.Message)
		return
	}
	response.WithCursorPagination(c, http.StatusOK, activities, domain.CursorPagination[time.Time]{
		Cursor:     cursor,
		NextCursor: nextCursor,
	})
}

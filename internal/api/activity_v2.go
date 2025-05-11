package api

import (
	"net/http"
	"time"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/services/activity_v2"
	"github.com/Stuhub-io/internal/api/decorators"
	"github.com/Stuhub-io/internal/api/middleware"
	"github.com/Stuhub-io/internal/api/request"
	"github.com/Stuhub-io/internal/api/response"
	"github.com/Stuhub-io/utils/pageutils"
	"github.com/gin-gonic/gin"
)

type ActivityV2Handler struct {
	activityV2Service *activity_v2.Service
	authMiddleware    *middleware.AuthMiddleware
}

type NewActivityV2HandlerParams struct {
	Router            *gin.RouterGroup
	ActivityV2Service *activity_v2.Service
	AuthMiddleware    *middleware.AuthMiddleware
}

func UseActivityV2Handler(params NewActivityV2HandlerParams) {
	handler := &ActivityV2Handler{
		activityV2Service: params.ActivityV2Service,
		authMiddleware:    params.AuthMiddleware,
	}

	router := params.Router.Group("/activity-services")
	authMiddleware := params.AuthMiddleware

	router.Use(authMiddleware.Authenticated())
	router.GET("/pages/:"+pageutils.PagePkIDParam+"/activities", decorators.RequiredAuth(decorators.CurrentUser(handler.ListActivitiesV2)))
}

func (h *ActivityV2Handler) ListActivitiesV2(c *gin.Context, curUser *domain.User) {
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

	// pagination based on activities endtime
	cursor := time.Now()
	if query.EndTime != "" {
		time, err := time.Parse(time.RFC3339, query.EndTime)
		if err != nil {
			response.BindError(c, err.Error())
			return
		}
		cursor = time
	}

	activities, err := h.activityV2Service.ListPageActivities(curUser, pagePkID, domain.CursorPagination[time.Time]{
		Limit:  query.Limit,
		Cursor: cursor,
	})

	nextCursorStr := domain.CalculateNextCursor[domain.ActivityV2, string](1, activities, "CreatedAt")
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

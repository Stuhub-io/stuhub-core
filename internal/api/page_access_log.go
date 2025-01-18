package api

import (
	"time"

	"github.com/Stuhub-io/core/domain"
	pageAccessLog "github.com/Stuhub-io/core/services/page_access_log"
	"github.com/Stuhub-io/internal/api/decorators"
	"github.com/Stuhub-io/internal/api/middleware"
	"github.com/Stuhub-io/internal/api/response"
	"github.com/gin-gonic/gin"
)

type PageAccessLogHandler struct {
	PageAccessLogService *pageAccessLog.Service
	AuthMiddleware       *middleware.AuthMiddleware
}

type NewPageAccessLogHandlerParams struct {
	Router               *gin.RouterGroup
	AuthMiddleware       *middleware.AuthMiddleware
	PageAccessLogService *pageAccessLog.Service
}

func UsePageAccessLogHandler(params NewPageAccessLogHandlerParams) {
	handler := &PageAccessLogHandler{
		PageAccessLogService: params.PageAccessLogService,
		AuthMiddleware:       params.AuthMiddleware,
	}
	router := params.Router.Group("/page-access-log-services")
	authMiddleware := params.AuthMiddleware

	router.Use(authMiddleware.Authenticated())
	router.GET("/logs", decorators.CurrentUser(handler.GetLogsList))
}

func (h *PageAccessLogHandler) GetLogsList(c *gin.Context, user *domain.User) {
	var queryParams struct {
		Cursor *time.Time `binding:"omitempty" form:"cursor"`
		Limit  int        `binding:"omitempty,gt=0"  form:"limit"`
	}

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var cursor time.Time
	if queryParams.Cursor.IsZero() {
		cursor = time.Now()
	} else {
		cursor = *queryParams.Cursor
	}

	if queryParams.Limit == 0 {
		queryParams.Limit = 20
	}

	cursorPagination := domain.CursorPagination[time.Time]{
		Cursor: cursor,
		Limit:  queryParams.Limit,
	}
	logs, nextCursor, err := h.PageAccessLogService.GetLogsByUser(cursorPagination, user)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithCursorPagination(c, 200, logs, domain.CursorPagination[*time.Time]{
		NextCursor: nextCursor,
		Limit:      queryParams.Limit,
	})
}

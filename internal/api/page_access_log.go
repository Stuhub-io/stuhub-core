package api

import (
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
	router.GET("/", decorators.CurrentUser(handler.GetLogsList))
}

func (h *PageAccessLogHandler) GetLogsList(c *gin.Context, user *domain.User) {
	var queryParams struct {
		Offset int `form:"offset" binding:"omitempty,gte=0"`
		Limit  int `form:"limit" binding:"omitempty,gt=0"`
	}

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if queryParams.Limit == 0 {
		queryParams.Limit = 20
	}

	logs, err := h.PageAccessLogService.GetLogsByUser(domain.OffsetBasedPagination{
		Offset: queryParams.Offset,
		Limit:  queryParams.Limit,
	}, user.PkID)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, 200, logs)
}

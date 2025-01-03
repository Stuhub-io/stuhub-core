package api

import (
	"context"

	"github.com/Stuhub-io/internal/api/middleware"
	"github.com/Stuhub-io/internal/api/response"
	"github.com/Stuhub-io/internal/repository/postgres"
	"github.com/gin-gonic/gin"
)

type PageAccessLogHandler struct {
	Repo           *postgres.PageAccessLogRepository
	AuthMiddleware *middleware.AuthMiddleware
}

type NewPageAccessLogHandlerParams struct {
	Router         *gin.RouterGroup
	AuthMiddleware *middleware.AuthMiddleware
	Repo           *postgres.PageAccessLogRepository
}

func UsePageAccessLogHandler(params NewPageAccessLogHandlerParams) {
	handler := &PageAccessLogHandler{
		Repo:           params.Repo,
		AuthMiddleware: params.AuthMiddleware,
	}
	router := params.Router.Group("/page-access-log-services")
	// authMiddleware := params.AuthMiddleware
	router.GET("/", handler.GetLogs)
}

func (h *PageAccessLogHandler) GetLogs(c *gin.Context) {
	logs, err := h.Repo.GetByUserPKID(context.Background(), 1)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, 200, logs)
}

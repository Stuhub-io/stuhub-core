package api

import (
	"net/http"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/services/page"
	"github.com/Stuhub-io/internal/api/decorators"
	"github.com/Stuhub-io/internal/api/middleware"
	"github.com/Stuhub-io/internal/api/request"
	"github.com/Stuhub-io/internal/api/response"
	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	pageService    *page.Service
	AuthMiddleware *middleware.AuthMiddleware
}

type NewPageHandlerParams struct {
	Router         *gin.RouterGroup
	AuthMiddleware *middleware.AuthMiddleware
	PageService    *page.Service
}

func UsePageHanlder(params NewPageHandlerParams) {
	handler := &PageHandler{
		pageService:    params.PageService,
		AuthMiddleware: params.AuthMiddleware,
	}

	router := params.Router.Group("/page-services")
	authMiddleware := params.AuthMiddleware

	router.Use(authMiddleware.Authenticated())
	router.POST("/create", decorators.CurrentUser(handler.CreateNewPage))
	router.GET("/all", decorators.CurrentUser(handler.GetSpacePages))
}

func (h *PageHandler) CreateNewPage(c *gin.Context, user *domain.User) {
	var body request.CreatePageBody

	if ok, verr := request.Validate(c, &body); !ok {
		response.BindError(c, verr.Error())
		return
	}

	page, err := h.pageService.CreateNewPage(page.CreatePageDto{
		ParentPagePkID: body.ParentPagePkID,
		SpacePkID:      body.SpacePkID,
		Name:           body.Name,
		ViewType:       domain.PageViewFromString(body.ViewType),
	})

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, http.StatusOK, page)
}

func (h *PageHandler) GetSpacePages(c *gin.Context, user *domain.User) {
	var params request.GetPagesBySpacePkIDParams
	if ok, verr := request.Validate(c, &params); !ok {
		response.BindError(c, verr.Error())
		return
	}
	pages, err := h.pageService.GetPagesBySpacePkID(params.SpacePkID)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}
	response.WithData(c, http.StatusOK, pages)
}

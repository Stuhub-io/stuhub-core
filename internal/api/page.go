package api

import (
	"net/http"
	"strconv"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/services/page"
	"github.com/Stuhub-io/internal/api/decorators"
	"github.com/Stuhub-io/internal/api/middleware"
	"github.com/Stuhub-io/internal/api/request"
	"github.com/Stuhub-io/internal/api/response"
	"github.com/Stuhub-io/utils/pageutils"
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
	router.POST("/pages", decorators.CurrentUser(handler.CreateNewPage))
	router.GET("/pages", decorators.CurrentUser(handler.GetSpacePages))
	router.GET("/pages/:"+pageutils.PageIDParam, decorators.CurrentUser(handler.GetPageByID))
	router.DELETE("/pages", decorators.CurrentUser(handler.DeletePageByPkID))
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

func (h *PageHandler) DeletePageByPkID(c *gin.Context, user *domain.User) {
	PagePkID := c.Query("id")
	if PagePkID == "" {
		response.WithErrorMessage(c, 400, "", "PkID Parameter Is Required")
		return
	}
	pkID, err := strconv.ParseInt(PagePkID, 10, 64)
	if err != nil {
		response.WithErrorMessage(c, 400, "", "Invalid PkID")
		return
	}
	data, domainErr := h.pageService.DeletePageByPkID(pkID, user.PkID)
	if domainErr != nil {
		response.WithErrorMessage(c, domainErr.Code, domainErr.Error, domainErr.Message)
		return
	}
	response.WithData(c, http.StatusOK, data)
}

func (h *PageHandler) GetPageByID(c *gin.Context, user *domain.User) {
	pageID, valid := pageutils.GetPageIDParam(c)
	if !valid {
		response.BindError(c, "pagePkID is missing or invalid")
		return
	}
	page, err := h.pageService.GetPageByID(pageID)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}
	response.WithData(c, http.StatusOK, page)
}

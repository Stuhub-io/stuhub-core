package api

import (
	"net/http"
	"path"
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
	router.GET("/pages", decorators.CurrentUser(handler.GetSpacePages))
	router.GET(path.Join("pages", ":"+pageutils.PageIDParam), decorators.CurrentUser(handler.GetPageByID))
	router.PUT(path.Join("pages", ":"+pageutils.PageIDParam), decorators.CurrentUser(handler.UpdatePageByID))
	router.POST(path.Join("pages", ":"+pageutils.PageIDParam, "archive"), decorators.CurrentUser(handler.ArchivedPageByID))

	router.DELETE("/pages", decorators.CurrentUser(handler.DeletePageByPkID))
	router.POST("/pages/bulk", decorators.CurrentUser(handler.BulkGetOrCreateByNodeID))
	router.DELETE("/pages/bulk", decorators.CurrentUser(handler.BulkArchivePages))
	// Depcrecated
	router.POST("/pages", decorators.CurrentUser(handler.CreateNewPage))
}

// Deprecated: create page with DocumentServices
func (h *PageHandler) CreateNewPage(c *gin.Context, user *domain.User) {
	var body request.CreatePageBody

	if ok, verr := request.Validate(c, &body); !ok {
		response.BindError(c, verr.Error())
		return
	}

	page, err := h.pageService.CreateNewPage(domain.PageInput{
		ParentPagePkID: body.ParentPagePkID,
		SpacePkID:      body.SpacePkID,
		Name:           body.Name,
		ViewType:       body.ViewType,
		CoverImage:     body.CoverImage,
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

func (h *PageHandler) UpdatePageByID(c *gin.Context, user *domain.User) {
	pageID, valid := pageutils.GetPageIDParam(c)
	if !valid {
		response.BindError(c, "pagePkID is missing or invalid")
		return
	}
	var params request.UpdatePageBody
	if ok, err := request.Validate(c, &params); !ok {
		response.BindError(c, err.Error())
		return
	}

	page, err := h.pageService.UpdatePageById(pageID, domain.PageInput{
		Name:           params.Name,
		ViewType:       params.ViewType,
		ParentPagePkID: params.ParentPagePkID,
		CoverImage:     params.CoverImage,
	})
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}
	response.WithData(c, http.StatusOK, page)
}

func (h *PageHandler) ArchivedPageByID(c *gin.Context, user *domain.User) {
	pageID, valid := pageutils.GetPageIDParam(c)
	if !valid {
		response.BindError(c, "pageID is missing or invalid")
		return
	}
	page, err := h.pageService.ArchivedPageByID(pageID)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}
	response.WithData(c, http.StatusOK, page)
}

func (h *PageHandler) BulkGetOrCreateByNodeID(c *gin.Context, user *domain.User) {
	var body request.BulkGetOrCreateByNodeIDBody

	if ok, verr := request.Validate(c, &body); !ok {
		response.BindError(c, verr.Error())
		return
	}

	inputs := make([]domain.PageInput, len(body.PageInputs))
	for i, page := range body.PageInputs {
		inputs[i] = domain.PageInput{
			Name:           page.Name,
			SpacePkID:      page.SpacePkID,
			ParentPagePkID: page.ParentPagePkID,
			ViewType:       page.ViewType,
			NodeID:         page.NodeID,
			CoverImage:     page.CoverImage,
		}
	}
	pages, err := h.pageService.BulkGetOrCreatePageByNodeID(inputs)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}
	response.WithData(c, http.StatusOK, pages)
}

func (h *PageHandler) BulkArchivePages(c *gin.Context, user *domain.User) {
	var body request.BulkArchivePagesBody

	if ok, verr := request.Validate(c, &body); !ok {
		response.BindError(c, verr.Error())
		return
	}
	if len(body.PagePkIDs) == 0 {
		response.BindError(c, "PagePkIDs is required")
		return
	}

	err := h.pageService.BulkArchivePages(body.PagePkIDs)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}
	response.WithData(c, http.StatusOK, nil)
}

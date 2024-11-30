package api

import (
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

func UsePageHandle(params NewPageHandlerParams) {
	handler := &PageHandler{
		pageService:    params.PageService,
		AuthMiddleware: params.AuthMiddleware,
	}
	router := params.Router.Group("/page-services")
	authMiddleware := params.AuthMiddleware

	router.Use(authMiddleware.Authenticated())
	router.GET("/pages", decorators.CurrentUser(handler.GetPages))
	router.POST("/pages", decorators.CurrentUser(handler.CreateDocument))
	router.GET("/pages/id/:"+pageutils.PageIDParam, decorators.CurrentUser(handler.GetPage))
	router.PUT(("/pages/:" + pageutils.PagePkIDParam), decorators.CurrentUser(handler.UpdatePage))
	router.PUT("/pages/:"+pageutils.PagePkIDParam+"/content", decorators.CurrentUser(handler.UpdatePageContent))
	router.PUT("/pages/:"+pageutils.PagePkIDParam+"/move", decorators.CurrentUser(handler.MovePage))
	router.DELETE("/pages/:"+pageutils.PagePkIDParam, decorators.CurrentUser(handler.ArchivePage))

	// asssets
	router.POST("pages/assets", decorators.CurrentUser(handler.CreateAsset))
}

func (h *PageHandler) GetPage(c *gin.Context, user *domain.User) {
	pageID, ok := pageutils.GetPageIDParam(c)
	if !ok {
		response.BindError(c, "pageID is missing or invalid")
		return
	}

	page, err := h.pageService.GetPageDetailByID(pageID)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, 200, page)
}

func (h *PageHandler) GetPages(c *gin.Context, user *domain.User) {
	var query request.GetPagesQuery
	if ok, verr := request.Validate(c, &query); !ok {
		response.BindError(c, verr.Error())
		return
	}

	pages, err := h.pageService.GetPagesByOrgPkID(domain.PageListQuery{
		OrgPkID:        query.OrgPkID,
		ViewTypes:      query.ViewTypes,
		ParentPagePkID: query.ParentPagePkID,
		Offset:         int(query.PaginationRequest.Page * query.PaginationRequest.Size),
		Limit:          int(query.PaginationRequest.Size),
		IsArchived:     query.IsArchived,
	})

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithPagination(c, 200, pages, domain.Pagination{
		Page: query.PaginationRequest.Page,
		Size: int64(len(pages)),
	})
}

func (h *PageHandler) CreateDocument(c *gin.Context, user *domain.User) {
	var body request.CreateDocumentBody

	if ok, verr := request.Validate(c, &body); !ok {
		response.BindError(c, verr.Error())
		return
	}

	if body.Document.JsonContent == "" {
		body.Document.JsonContent = "{}"
	}

	page, err := h.pageService.CreateDocumentPage(domain.DocumentPageInput{
		PageInput: domain.PageInput{
			Name:             body.Name,
			ParentPagePkID:   body.ParentPagePkID,
			ViewType:         body.ViewType,
			CoverImage:       body.CoverImage,
			OrganizationPkID: body.OrgPkID,
		},
		Document: domain.DocumentInput(body.Document),
	})

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, 200, page)
}

func (h *PageHandler) UpdatePage(c *gin.Context, user *domain.User) {
	pagePkID, ok := pageutils.GetPagePkIDParam(c)
	if !ok {
		response.BindError(c, "pagePkID is missing or invalid")
		return
	}

	var body request.UpdatePageBody
	if ok, verr := request.Validate(c, &body); !ok {
		response.BindError(c, verr.Error())
		return
	}

	if body.Document != nil && body.Document.JsonContent == "" {
		body.Document.JsonContent = "{}"
	}

	document, err := h.pageService.UpdatePageByPkID(pagePkID, domain.PageUpdateInput{
		Name:       body.Name,
		ViewType:   body.ViewType,
		CoverImage: body.CoverImage,
		Document:   body.Document,
	})
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}
	response.WithData(c, 200, document)
}

func (h *PageHandler) UpdatePageContent(c *gin.Context, user *domain.User) {
	pagePkID, ok := pageutils.GetPagePkIDParam(c)

	if !ok {
		response.BindError(c, "pagePkID is missing or invalid")
	}
	var body request.UpdatePageContent
	if ok, verr := request.Validate(c, &body); !ok {
		response.BindError(c, verr.Error())
		return
	}

	document, err := h.pageService.UpdateDocumentContentByPkID(pagePkID, domain.DocumentInput{
		JsonContent: body.JsonContent,
	})
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}
	response.WithData(c, 200, document)
}

func (h *PageHandler) ArchivePage(c *gin.Context, user *domain.User) {
	pagePkID, ok := pageutils.GetPagePkIDParam(c)
	if !ok {
		response.BindError(c, "pagePkID is missing or invalid")
	}
	page, err := h.pageService.ArchivedPageByPkID(pagePkID)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}
	response.WithData(c, 200, page)
}

func (h *PageHandler) MovePage(c *gin.Context, user *domain.User) {
	pagePkID, ok := pageutils.GetPagePkIDParam(c)
	if !ok {
		response.BindError(c, "pagePkID is missing or invalid")
	}

	var body request.MovePageBody
	if ok, verr := request.Validate(c, &body); !ok {
		response.BindError(c, verr.Error())
		return
	}

	page, err := h.pageService.MovePageByPkID(pagePkID, domain.PageMoveInput{
		ParentPagePkID: body.ParentPagePkID,
	})
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}
	response.WithData(c, 200, page)
}

// Assets

func (h *PageHandler) CreateAsset(c *gin.Context, user *domain.User) {
	var body request.CreateAssetBody

	if ok, verr := request.Validate(c, &body); !ok {
		response.BindError(c, verr.Error())
		return
	}

	page, err := h.pageService.CreateAssetPage(domain.AssetPageInput{
		PageInput: domain.PageInput{
			Name:             body.Name,
			ParentPagePkID:   body.ParentPagePkID,
			ViewType:         body.ViewType,
			CoverImage:       body.CoverImage,
			OrganizationPkID: body.OrgPkID,
		},
		Asset: domain.AssetInput{
			URL:        body.Asset.Url,
			Size:       body.Asset.Size,
			Extension:  body.Asset.Extension,
			Thumbnails: body.Asset.Thumbnails,
		},
	})

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, 200, page)
}

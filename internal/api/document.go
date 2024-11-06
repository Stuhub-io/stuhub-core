package api

import (
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/services/document"
	"github.com/Stuhub-io/internal/api/decorators"
	"github.com/Stuhub-io/internal/api/middleware"
	"github.com/Stuhub-io/internal/api/request"
	"github.com/Stuhub-io/internal/api/response"
	"github.com/Stuhub-io/utils/docutils"
	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	documentService *document.Service
	AuthMiddleware  *middleware.AuthMiddleware
}

type NewDocumentHandlerParams struct {
	Router          *gin.RouterGroup
	AuthMiddleware  *middleware.AuthMiddleware
	DocumentService *document.Service
}

func UseDocumentHandle(params NewDocumentHandlerParams) {
	handler := &DocumentHandler{
		documentService: params.DocumentService,
		AuthMiddleware:  params.AuthMiddleware,
	}
	router := params.Router.Group("/document-services")
	authMiddleware := params.AuthMiddleware

	router.Use(authMiddleware.Authenticated())
	router.GET("/documents", decorators.CurrentUser(handler.GetDocuments))
	router.POST("/documents", decorators.CurrentUser(handler.CreateDocument))
	router.GET("/documents/id/:"+docutils.PageIDParam, decorators.CurrentUser(handler.GetDocument))
	router.PUT(("/documents/:" + docutils.PagePkIDParam), decorators.CurrentUser(handler.UpdateDocument))
	router.DELETE("/documents"+docutils.PagePkIDParam, decorators.CurrentUser(handler.ArchiveDocument))
}

func (h *DocumentHandler) GetDocument(c *gin.Context, user *domain.User) {
	pageID, ok := docutils.GetPageIDParam(c)
	if !ok {
		response.BindError(c, "pageID is missing or invalid")
		return
	}

	page, err := h.documentService.GetPageDetailByID(pageID)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, 200, page)
}

func (h *DocumentHandler) GetDocuments(c *gin.Context, user *domain.User) {
	var query request.GetPagesQuery
	if ok, verr := request.Validate(c, &query); !ok {
		response.BindError(c, verr.Error())
		return
	}

	pages, err := h.documentService.GetPagesByOrgPkID(domain.PageListQuery{
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

func (h *DocumentHandler) CreateDocument(c *gin.Context, user *domain.User) {
	var body request.CreatePageBody

	if ok, verr := request.Validate(c, &body); !ok {
		response.BindError(c, verr.Error())
		return
	}

	if body.Document.JsonContent == "" {
		body.Document.JsonContent = "{}"
	}

	page, err := h.documentService.CreatePage(domain.PageInput{
		Name:             body.Name,
		ParentPagePkID:   body.ParentPagePkID,
		ViewType:         body.ViewType,
		CoverImage:       body.CoverImage,
		OrganizationPkID: body.OrgPkID,
		Document:         domain.DocumentInput(body.Document),
	})

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, 200, page)
}

func (h *DocumentHandler) UpdateDocument(c *gin.Context, user *domain.User) {
	pagePkID, ok := docutils.GetPagePkIDParam(c)
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

	document, err := h.documentService.UpdatePageByPkID(pagePkID, domain.PageUpdateInput{
		Name:           body.Name,
		ParentPagePkID: body.ParentPagePkID,
		ViewType:       body.ViewType,
		CoverImage:     body.CoverImage,
		Document:       body.Document,
	})
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}
	response.WithData(c, 200, document)
}

func (h *DocumentHandler) ArchiveDocument(c *gin.Context, user *domain.User) {
}

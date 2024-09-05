package api

import (
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/services/document"
	"github.com/Stuhub-io/core/services/page"
	"github.com/Stuhub-io/internal/api/decorators"
	"github.com/Stuhub-io/internal/api/middleware"
	"github.com/Stuhub-io/internal/api/request"
	"github.com/Stuhub-io/internal/api/response"
	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	documentService *document.Service
	pageService     *page.Service
	AuthMiddleware  *middleware.AuthMiddleware
}

type NewDocumentHandlerParams struct {
	Router          *gin.RouterGroup
	AuthMiddleware  *middleware.AuthMiddleware
	DocumentService *document.Service
	PageService     *page.Service
}

func UseDocumentHandle(params NewDocumentHandlerParams) {
	handler := &DocumentHandler{
		documentService: params.DocumentService,
		pageService:     params.PageService,
		AuthMiddleware:  params.AuthMiddleware,
	}
	router := params.Router.Group("/document-services")
	authMiddleware := params.AuthMiddleware
	router.Use(authMiddleware.Authenticated())
	router.POST("/create", decorators.CurrentUser(handler.CreateNewDocument))
	// router.GET("/get/:documentPkID", decorators.CurrentUser(handler))
}

func (h *DocumentHandler) CreateNewDocument(c *gin.Context, user *domain.User) {
	var body request.CreateDocumentBody

	if ok, verr := request.Validate(c, &body); !ok {
		response.BindError(c, verr.Error())
		return
	}

	pageInput := body.Page

	page, err := h.pageService.CreateNewPage(page.CreatePageDto{
		Name:           pageInput.Name,
		SpacePkID:      pageInput.SpacePkID,
		ParentPagePkID: pageInput.ParentPagePkID,
		ViewType:       domain.PageViewFromString(pageInput.ViewType),
	})
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	doc, err := h.documentService.CreateNewDocument(
		page.PkId,
		body.JsonContent,
	)

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}
	type DocumentResponse struct {
		Page domain.Page     `json:"page"`
		Doc  domain.Document `json:"document"`
	}

	response.WithData(c, 200, DocumentResponse{
		Page: *page,
		Doc:  *doc,
	})
}

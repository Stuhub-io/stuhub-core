package api

import (
	"fmt"

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
	router.GET(("/pages/quick-search"), decorators.CurrentUser(handler.QuickSearch))

	router.PUT(
		"/pages/:"+pageutils.PagePkIDParam+"/content",
		decorators.CurrentUser(handler.UpdatePageContent),
	)
	router.PUT("/pages/:"+pageutils.PagePkIDParam+"/move", decorators.CurrentUser(handler.MovePage))
	router.DELETE("/pages/:"+pageutils.PagePkIDParam, decorators.CurrentUser(handler.ArchivePage))

	// public page
	router.POST(
		"pages/id/:"+pageutils.PageIDParam+"/public-token",
		decorators.CurrentUser(handler.CreatePagePublicToken),
	)
	router.DELETE(
		"pages/id/:"+pageutils.PageIDParam+"/public-token",
		decorators.CurrentUser(handler.ArchiveAllPagePublicToken),
	)
	router.GET("pages/public-token/:"+pageutils.PublicTokenIDParam, handler.GetPageByToken)

	// asssets
	router.POST("pages/assets", decorators.CurrentUser(handler.CreateAsset))

	// page roles
	router.PUT(
		("/pages/:" + pageutils.PagePkIDParam + "/general-access"),
		decorators.CurrentUser(handler.UpdatePageGeneralAccess),
	)
	router.POST(
		("/pages/:" + pageutils.PagePkIDParam + "/roles"),
		decorators.CurrentUser(handler.AddPageRoleUser),
	)
	router.GET(
		("/pages/:" + pageutils.PagePkIDParam + "/roles"),
		decorators.CurrentUser(handler.GetAllRoleUsers),
	)
	router.PATCH(
		("/pages/:" + pageutils.PagePkIDParam + "/roles"),
		decorators.CurrentUser(handler.UpdatePageRoleUser),
	)
	router.DELETE(
		("/pages/:" + pageutils.PagePkIDParam + "/roles"),
		decorators.CurrentUser(handler.DeletePageRoleUser),
	)

	// page role requests
	router.POST(
		("/pages/id/:" + pageutils.PageIDParam + "/role-requests"),
		decorators.RequiredAuth(decorators.CurrentUser(handler.RequestPageAccess)),
	)

	router.GET(
		("/pages/:" + pageutils.PagePkIDParam + "/role-requests"),
		decorators.RequiredAuth(decorators.CurrentUser(handler.ListRequestPageAccesses)),
	)
	router.POST(
		("/pages/:" + pageutils.PagePkIDParam + "/role-requests/accept"),
		decorators.RequiredAuth(decorators.CurrentUser(handler.AcceptRequestPageAccess)),
	)
	router.POST(
		("/pages/:" + pageutils.PagePkIDParam + "/role-requests/reject"),
		decorators.RequiredAuth(decorators.CurrentUser(handler.RejectRequestPageAccesses)),
	)
	router.POST(
		("/pages/:" + pageutils.PagePkIDParam + "/star"),
		decorators.RequiredAuth(decorators.CurrentUser(handler.AddPageToStarred)),
	)
	router.POST(
		("/pages/:" + pageutils.PagePkIDParam + "/unstar"),
		decorators.RequiredAuth(decorators.CurrentUser(handler.RemovePageFromStarred)),
	)
}

func (h *PageHandler) GetPage(c *gin.Context, user *domain.User) {
	pageID, ok := pageutils.GetPageIDParam(c)

	if !ok {
		response.BindError(c, "pageID is missing or invalid")
		return
	}

	page, err := h.pageService.GetPageDetailByID(pageID, "", user)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, 200, page)
}

func (h *PageHandler) GetPages(c *gin.Context, user *domain.User) {
	var query request.GetPagesQuery
	if verr := request.Validate(c, &query); verr != nil {
		response.BindError(c, verr.Error())
		return
	}

	// get starred page if only authenticated
	var starredByUserPkID *int64 = nil
	if query.IsStarred != nil && *query.IsStarred && user != nil {
		starredByUserPkID = &user.PkID
	}

	pages, err := h.pageService.GetPagesByOrgPkID(domain.PageListQuery{
		OrgPkID:             &query.OrgPkID,
		ViewTypes:           query.ViewTypes,
		ParentPagePkID:      query.ParentPagePkID,
		Offset:              int(query.PaginationRequest.Page * query.PaginationRequest.Size),
		IsAll:               query.All,
		Limit:               int(query.PaginationRequest.Size),
		IsArchived:          query.IsArchived,
		GeneralRole:         query.GeneralRole,
		IsStarredByUserPkID: starredByUserPkID,
	}, user)

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

	if verr := request.Validate(c, &body); verr != nil {
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
			AuthorPkID:       user.PkID,
			OrganizationPkID: body.OrgPkID,
		},
		Document: domain.DocumentInput(body.Document),
	}, user)

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
	if verr := request.Validate(c, &body); verr != nil {
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
	}, user)

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
	if verr := request.Validate(c, &body); verr != nil {
		response.BindError(c, verr.Error())
		return
	}

	document, err := h.pageService.UpdateDocumentContentByPkID(pagePkID, domain.DocumentInput{
		JsonContent: body.JsonContent,
	}, user)
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
	page, err := h.pageService.ArchivedPageByPkID(pagePkID, user)
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
	if verr := request.Validate(c, &body); verr != nil {
		response.BindError(c, verr.Error())
		return
	}

	page, err := h.pageService.MovePageByPkID(pagePkID, domain.PageMoveInput{
		ParentPagePkID: body.ParentPagePkID,
	}, user)

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}
	response.WithData(c, 200, page)
}

// Assets

func (h *PageHandler) CreateAsset(c *gin.Context, user *domain.User) {
	var body request.CreateAssetBody

	if verr := request.Validate(c, &body); verr != nil {
		response.BindError(c, verr.Error())
		return
	}

	page, err := h.pageService.CreateAssetPage(domain.AssetPageInput{
		PageInput: domain.PageInput{
			Name:             body.Name,
			ParentPagePkID:   body.ParentPagePkID,
			ViewType:         body.ViewType,
			CoverImage:       body.CoverImage,
			AuthorPkID:       user.PkID,
			OrganizationPkID: body.OrgPkID,
		},
		Asset: domain.AssetInput{
			URL:        body.Asset.Url,
			Size:       body.Asset.Size,
			Extension:  body.Asset.Extension,
			Thumbnails: body.Asset.Thumbnails,
		},
	}, user)

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, 200, page)
}

// public Page.
func (h *PageHandler) CreatePagePublicToken(c *gin.Context, user *domain.User) {
	pageID, ok := pageutils.GetPageIDParam(c)
	if !ok {
		response.BindError(c, "pageID is missing or invalid")
		return
	}

	token, err := h.pageService.CreatePublicPageToken(pageID)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, 200, token)
}

func (h *PageHandler) ArchiveAllPagePublicToken(c *gin.Context, user *domain.User) {
	pageID, ok := pageutils.GetPageIDParam(c)
	if !ok {
		response.BindError(c, "pageID is missing or invalid")
		return
	}

	err := h.pageService.ArchiveAllPublicPageToken(pageID)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, 200, nil)
}

func (h *PageHandler) UpdatePageGeneralAccess(c *gin.Context, user *domain.User) {
	pagePkID, ok := pageutils.GetPagePkIDParam(c)
	if !ok {
		response.BindError(c, "pagePkID is missing or invalid")
		return
	}

	var body request.UpdatePageGeneralAccessBody
	if verr := request.Validate(c, &body); verr != nil {
		response.BindError(c, verr.Error())
		return
	}

	page, err := h.pageService.UpdateGeneralAccess(pagePkID, domain.PageGeneralAccessUpdateInput{
		AuthorPkID:  user.PkID,
		GeneralRole: body.GeneralRole,
	}, user)

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, 200, page)
}

func (h *PageHandler) GetPageByToken(c *gin.Context) {
	tokenID, ok := pageutils.GetPublicTokenIDParam(c)
	if !ok {
		response.BindError(c, "token is missing or invalid")
		return
	}

	page, err := h.pageService.GetPageDetailByID("", tokenID, nil)

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, 200, page)
}

// Page Roles.
func (h *PageHandler) AddPageRoleUser(c *gin.Context, user *domain.User) {
	pagePkID, ok := pageutils.GetPagePkIDParam(c)
	if !ok {
		response.BindError(c, "pagePkID is missing or invalid")
		return
	}

	var body request.AddPageRoleUserBody
	if verr := request.Validate(c, &body); verr != nil {
		response.BindError(c, verr.Error())
		return
	}

	pageRole, _, err := h.pageService.AddPageRoleUser(domain.PageRoleCreateInput{
		PagePkID: pagePkID,
		Email:    body.Email,
		Role:     body.Role,
	}, user)

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, 201, pageRole)
}

func (h *PageHandler) GetAllRoleUsers(c *gin.Context, user *domain.User) {
	pagePkID, ok := pageutils.GetPagePkIDParam(c)
	if !ok {
		response.BindError(c, "pagePkID is missing or invalid")
		return
	}

	roles, err := h.pageService.GetPageRoleUsers(domain.PageRoleGetAllInput{
		AuthorPkID: user.PkID,
		PagePkID:   pagePkID,
	}, user)

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, 200, roles)
}

func (h *PageHandler) UpdatePageRoleUser(c *gin.Context, user *domain.User) {
	pagePkID, ok := pageutils.GetPagePkIDParam(c)
	if !ok {
		response.BindError(c, "pagePkID is missing or invalid")
		return
	}

	var body request.UpdatePageRoleUserBody
	if verr := request.Validate(c, &body); verr != nil {
		response.BindError(c, verr.Error())
		return
	}

	err := h.pageService.UpdatePageRoleUser(domain.PageRoleUpdateInput{
		AuthorPkID: user.PkID,
		PagePkID:   pagePkID,
		Email:      body.Email,
		Role:       body.Role,
	}, user)

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithMessage(c, 200, "Role's updated successfully!")
}

func (h *PageHandler) DeletePageRoleUser(c *gin.Context, user *domain.User) {
	pagePkID, ok := pageutils.GetPagePkIDParam(c)
	if !ok {
		response.BindError(c, "pagePkID is missing or invalid")
		return
	}

	var body request.DeletePageRoleUserBody
	if verr := request.Validate(c, &body); verr != nil {
		response.BindError(c, verr.Error())
		return
	}

	err := h.pageService.DeletePageRoleUser(domain.PageRoleDeleteInput{
		AuthorPkID: user.PkID,
		PagePkID:   pagePkID,
		Email:      body.Email,
	}, user)

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithMessage(c, 200, "Role's deleted successfully!")
}

func (h *PageHandler) RequestPageAccess(c *gin.Context, user *domain.User) {
	pageID, ok := pageutils.GetPageIDParam(c)
	if !ok {
		response.BindError(c, "pagePkID is missing or invalid")
		return
	}

	er := h.pageService.RequestPagePermission(pageID, user.Email)

	if er != nil {
		response.WithErrorMessage(c, er.Code, er.Error, er.Message)
		return
	}

	response.WithMessage(c, 200, "Request sent successfully!")
}

func (h *PageHandler) ListRequestPageAccesses(c *gin.Context, user *domain.User) {
	pagePkID, ok := pageutils.GetPagePkIDParam(c)
	if !ok {
		response.BindError(c, "pagePkID is missing or invalid")
		return
	}

	requests, err := h.pageService.ListRequestPagePermissions(pagePkID)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, 200, requests)
}

func (h *PageHandler) AcceptRequestPageAccess(c *gin.Context, user *domain.User) {
	pagePkID, ok := pageutils.GetPagePkIDParam(c)
	if !ok {
		response.BindError(c, "pagePkID is missing or invalid")
		return
	}

	var body request.AcceptRequestPageAccess
	if verr := request.Validate(c, &body); verr != nil {
		response.BindError(c, verr.Error())
		return
	}
	err := h.pageService.AcceptRequestPagePermission(domain.PageRoleCreateInput{
		PagePkID: pagePkID,
		Email:    body.Email,
		Role:     body.Role,
	}, user)

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithMessage(c, 200, "Request accepted successfully!")
}

func (h *PageHandler) RejectRequestPageAccesses(c *gin.Context, user *domain.User) {
	pagePkID, ok := pageutils.GetPagePkIDParam(c)
	if !ok {
		response.BindError(c, "pagePkID is missing or invalid")
		return
	}

	var body request.RejectRequestPageAccess
	if verr := request.Validate(c, &body); verr != nil {
		response.BindError(c, verr.Error())
		return
	}

	err := h.pageService.RejectPagePermissions(pagePkID, body.Emails)

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithMessage(c, 200, "Request rejected successfully!")
}

func (h *PageHandler) AddPageToStarred(c *gin.Context, user *domain.User) {
	pagePkID, ok := pageutils.GetPagePkIDParam(c)
	if !ok {
		response.BindError(c, "pagePkID is missing or invalid")
		return
	}

	var body request.ToggleStarPageBody
	if verr := request.Validate(c, &body); verr != nil {
		response.BindError(c, verr.Error())
		return
	}

	err := h.pageService.AddPageToStarred(domain.StarPageInput{
		PagePkID:      pagePkID,
		ActorUserPkID: user.PkID,
	}, user)

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithMessage(c, 200, "Page added to starred successfully!")
}

func (h *PageHandler) RemovePageFromStarred(c *gin.Context, user *domain.User) {
	pagePkID, ok := pageutils.GetPagePkIDParam(c)
	if !ok {
		response.BindError(c, "pagePkID is missing or invalid")
		return
	}

	var body request.ToggleStarPageBody
	if verr := request.Validate(c, &body); verr != nil {
		response.BindError(c, verr.Error())
		return
	}

	err := h.pageService.RemovePageFromStarred(domain.StarPageInput{
		PagePkID:      pagePkID,
		ActorUserPkID: user.PkID,
	}, user)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithMessage(c, 200, "Page added to starred successfully!")
}

func (h *PageHandler) QuickSearch(c *gin.Context, user *domain.User) {
	var queryParams struct {
		Keyword    string  `binding:"required" form:"keyword"`
		ViewType   *string `binding:"omitempty" form:"view_type"`
		AuthorPkID *int64  `binding:"omitempty" form:"author_pkid"`
	}

	if err := c.ShouldBindQuery(&queryParams); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	fmt.Print(">>> Params: ", queryParams)

	pages, err := h.pageService.QuickSearch(domain.SearchIndexedPageParams{
		UserPkID:   user.PkID,
		Keyword:    queryParams.Keyword,
		ViewType:   queryParams.ViewType,
		AuthorPkID: queryParams.AuthorPkID,
	})
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, 200, pages)
}

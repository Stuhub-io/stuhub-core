package api

import (
	"net/http"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/services/organization"
	"github.com/Stuhub-io/internal/api/decorators"
	"github.com/Stuhub-io/internal/api/middleware"
	"github.com/Stuhub-io/internal/api/request"
	"github.com/Stuhub-io/internal/api/response"
	"github.com/gin-gonic/gin"
)

type OrganizationHandler struct {
	orgService *organization.Service
}

type NewOrganizationHandlerParams struct {
	Router         *gin.RouterGroup
	AuthMiddleware *middleware.AuthMiddleware
	OrgService     *organization.Service
}

func UseOrganizationHandler(params NewOrganizationHandlerParams) {
	handler := &OrganizationHandler{
		orgService: params.OrgService,
	}

	router := params.Router.Group("/organization-services")
	authMiddleware := params.AuthMiddleware

	router.Use(authMiddleware.Authenticated())

	router.POST("/create", decorators.CurrentUser(handler.CreateOrganization))
	router.GET("/joined", decorators.CurrentUser(handler.GetJoinedOrgs))

}

func (h *OrganizationHandler) CreateOrganization(c *gin.Context, user *domain.User) {
	var body request.CreateOrgBody
	if ok, vr := request.Validate(c, &body); !ok {
		response.BindError(c, vr.Error())
		return
	}

	data, err := h.orgService.CreateOrganization(organization.CreateOrganizationParams{
		OwnerPkID:   user.PkID,
		Name:        body.Name,
		Description: body.Description,
		Avatar:      body.Avatar,
	})
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, http.StatusOK, data, "Success")
}

func (h *OrganizationHandler) GetJoinedOrgs(c *gin.Context, user *domain.User) {
	data, err := h.orgService.GetJoinedOrgs(user.PkID)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, http.StatusOK, data, "Success")
}
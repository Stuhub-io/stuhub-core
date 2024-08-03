package api

import (
	"net/http"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/services/space"
	"github.com/Stuhub-io/internal/api/decorators"
	"github.com/Stuhub-io/internal/api/middleware"
	"github.com/Stuhub-io/internal/api/request"
	"github.com/Stuhub-io/internal/api/response"
	"github.com/gin-gonic/gin"
)

type SpaceHandler struct {
	spaceService *space.Service
}

type NewSpaceHandlerParams struct {
	Router         *gin.RouterGroup
	AuthMiddleware *middleware.AuthMiddleware
	SpaceService   *space.Service
}

func UseSpaceHandler(params NewSpaceHandlerParams) {
	handler := &SpaceHandler{
		spaceService: params.SpaceService,
	}

	router := params.Router.Group("/space-services")
	authMiddleware := params.AuthMiddleware

	router.Use(authMiddleware.Authenticated())

	router.POST("/create", decorators.CurrentUser(handler.CreateNewSpace))
	router.GET("/joined", decorators.CurrentUser(handler.GetJoinSpaceByOrgPkID))
}

func (h *SpaceHandler) CreateNewSpace(c *gin.Context, user *domain.User) {
	var body request.CreateSpaceBody
	if ok, verr := request.Validate(c, &body); !ok {
		response.BindError(c, verr.Error())
		return
	}
	space, err := h.spaceService.CreateOrgSpace(space.CreateSpaceDto{
		OrgPkID:     body.OrgPkID,
		OwnerPkID:   user.PkID,
		Name:        body.Name,
		Description: body.Description,
	})
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}
	response.WithData(c, http.StatusOK, space)
}

func (h *SpaceHandler) GetJoinSpaceByOrgPkID(c *gin.Context, user *domain.User) {
	var params request.GetSpaceByOrgPkIDParams
	if ok, er := request.Validate(c, &params); !ok {
		response.BindError(c, er.Error())
		return
	}
	data, err := h.spaceService.GetJoinedSpaceByOrgPkID(params.OrgPkID, user.PkID)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}
	response.WithData(c, http.StatusOK, data)
}

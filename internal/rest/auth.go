package rest

import (
	"net/http"

	"github.com/Stuhub-io/core/services/auth"
	"github.com/Stuhub-io/internal/rest/request"
	"github.com/Stuhub-io/internal/rest/response"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService auth.Service
}

type NewAuthHandlerParams struct {
	Router      *gin.RouterGroup
	AuthService *auth.Service
}

func UseAuthHandler(params NewAuthHandlerParams) {
	handler := &AuthHandler{
		authService: *params.AuthService,
	}

	router := params.Router.Group("/auth-services")

	router.POST("/auth-email", handler.AuthenByEmail)
}

func (h *AuthHandler) AuthenByEmail(c *gin.Context) {
	var body request.RegisterByEmailBody

	if ok, vr := request.Validate(c, &body); !ok {
		response.BindError(c, vr.Error())
		return
	}

	data, err := h.authService.AuthenByEmailStepOne(auth.AuthenByEmailStepOneDto{
		Email: body.Email,
	})

	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, http.StatusOK, data, "Success")
}

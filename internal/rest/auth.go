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

	router.POST("/email-step-one", handler.AuthenByEmailStepOne)
	router.POST("/validate-email-token", handler.ValidateEmailToken)
	router.POST("/set-password", handler.SetPassword)
}

func (h *AuthHandler) AuthenByEmailStepOne(c *gin.Context) {
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

func (h *AuthHandler) ValidateEmailToken(c *gin.Context) {
	var body request.ValidateEmailTokenBody
	if ok, vr := request.Validate(c, &body); !ok {
		response.BindError(c, vr.Error())
		return
	}

	data, err := h.authService.ValidateEmailAuth(body.Token)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}
	response.WithData(c, http.StatusOK, data, "Success")
}

func (h *AuthHandler) SetPassword(c *gin.Context) {
	var body request.SetUserPasswordBody
	if ok, vr := request.Validate(c, &body); !ok {
		response.BindError(c, vr.Error())
		return
	}

	data, err := h.authService.SetPasswordAndAuthUser(auth.AuthenByEmailPassword{
		Email:       body.Email,
		Password:    body.Password,
		ActionToken: body.ActionToken,
	})
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}
	response.WithData(c, http.StatusOK, data, "Success")
}

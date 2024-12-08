package api

import (
	"net/http"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/services/auth"
	"github.com/Stuhub-io/internal/api/request"
	"github.com/Stuhub-io/internal/api/response"
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
	router.POST("/email", handler.AuthenUserByEmailPassword)
	router.POST("/google", handler.AuthenUserByGoogle)
	router.POST("/user-by-token", handler.GetUserByAccessToken)
}

func (h *AuthHandler) AuthenByEmailStepOne(c *gin.Context) {
	var body request.RegisterByEmailBody

	if vr := request.Validate(c, &body); vr != nil {
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
	if vr := request.Validate(c, &body); vr != nil {
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
	if vr := request.Validate(c, &body); vr != nil {
		response.BindError(c, vr.Error())
		return
	}

	data, err := h.authService.SetPasswordAndAuthUser(auth.AuthenByEmailAfterSetPasswordDto{
		Email:       body.Email,
		RawPassword: body.Password,
		ActionToken: body.ActionToken,
	})
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, http.StatusOK, data, "Success")
}

func (h *AuthHandler) AuthenUserByEmailPassword(c *gin.Context) {
	var body request.AuthenUserByEmailPasswordBody
	if vr := request.Validate(c, &body); vr != nil {
		response.BindError(c, vr.Error())
		return
	}
	token, user, err := h.authService.AuthenUserByEmailPassword(auth.AuthenByEmailPasswordDto{
		Email:       body.Email,
		RawPassword: body.Password,
	})
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	data := struct {
		Tokens  domain.AuthToken `json:"tokens"`
		Profile domain.User      `json:"profile"`
	}{
		Tokens:  *token,
		Profile: *user,
	}

	response.WithData(c, http.StatusOK, data, "Success")
}

func (h *AuthHandler) AuthenUserByGoogle(c *gin.Context) {
	var body request.AuthenUserByGoogleBody
	if vr := request.Validate(c, &body); vr != nil {
		response.BindError(c, vr.Error())
		return
	}

	data, err := h.authService.AuthenUserByGoogle(auth.AuthenByGoogleDto{
		Token: body.Token,
	})
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, http.StatusOK, data, "Success")
}

func (h *AuthHandler) GetUserByAccessToken(c *gin.Context) {
	var query request.GetUserByTokenQuery
	if vr := request.Validate(c, &query); vr != nil {
		response.BindError(c, vr.Error())
		return
	}

	user, err := h.authService.GetUserByToken(query.AccessToken)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, http.StatusOK, user, "Success")
}

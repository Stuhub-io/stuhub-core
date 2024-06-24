package rest

import (
	"net/http"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/services/auth"
	"github.com/Stuhub-io/internal/rest/request"
	"github.com/Stuhub-io/internal/rest/response"
	"github.com/gin-gonic/gin"
)

type AuthService interface {
	RegisterByEmail(loginDto auth.RegisterByEmailDto) *domain.Error
}

type AuthHandler struct {
	authService AuthService
}

type NewAuthHandlerParams struct {
	Router *gin.RouterGroup
	AuthService
}

func UseAuthHandler(params NewAuthHandlerParams) {
	handler := &AuthHandler{
		authService: params.AuthService,
	}

	router := params.Router.Group("/auth")

	router.POST("/register-email", handler.RegisterByEmail)
}

func (h *AuthHandler) RegisterByEmail(c *gin.Context) {
	var body request.RegisterByEmailBody

	if ok, vr := request.Validate(c, &body); !ok {
		response.BindError(c, vr.Error())
		return
	}

	err := h.authService.RegisterByEmail(auth.RegisterByEmailDto{
		Email: body.Email,
	})
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithMessage(c, http.StatusOK, "Check your email for verification")
}

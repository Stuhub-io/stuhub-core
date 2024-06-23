package rest

import (
	"net/http"

	"github.com/Stuhub-io/core/services/user"
	"github.com/Stuhub-io/internal/rest/response"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	Login(loginDto user.LoginDto) (*user.LoginResponse, error)
}

type UserHandler struct {
	userService UserService
}

type NewUserHandlerParams struct {
	Router *gin.Engine
	UserService
}

func UseUserHandler(params NewUserHandlerParams) {
	handler := &UserHandler{
		userService: params.UserService,
	}

	router := params.Router

	router.GET("/login", handler.Login)
	router.GET("/me", handler.Me)
}

func (h *UserHandler) Login(c *gin.Context) {
	loginDto := user.LoginDto{
		Username: "Khoa",
		Password: "123",
	}
	err := c.Bind(&loginDto)
	if err != nil {
		response.BindError(c, err.Error())
	}

	resp, err := h.userService.Login(loginDto)
	if err != nil {
		response.BadRequest(c)
	}

	response.WithData(c, http.StatusOK, resp)
}

func (h *UserHandler) Me(c *gin.Context) {
	loginDto := user.LoginDto{
		Username: "Khoa",
		Password: "123",
	}
	err := c.Bind(&loginDto)
	if err != nil {
		response.BindError(c, err.Error())
	}

	_, err = h.userService.Login(loginDto)
	if err != nil {
		c.JSON(getStatusCode(err), err.Error())
	}

	response.Unauthorized(c)
}

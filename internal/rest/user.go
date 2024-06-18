package rest

import (
	"net/http"

	"github.com/Stuhub-io/core/services/user"
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

func NewUserHandler(params NewUserHandlerParams) {
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
		c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	resp, err := h.userService.Login(loginDto)
	if err != nil {
		c.JSON(getStatusCode(err), err.Error())
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Me(c *gin.Context) {
	loginDto := user.LoginDto{
		Username: "Khoa",
		Password: "123",
	}
	err := c.Bind(&loginDto)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	resp, err := h.userService.Login(loginDto)
	if err != nil {
		c.JSON(getStatusCode(err), err.Error())
	}

	c.JSON(http.StatusOK, resp)
}

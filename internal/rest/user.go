package rest

import (
	"net/http"
	"strconv"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/services/user"
	"github.com/Stuhub-io/internal/rest/response"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	GetUserById(id int64) (*user.GetUserByIdResponse, *domain.Error)
}

type UserHandler struct {
	userService UserService
}

type NewUserHandlerParams struct {
	Router *gin.RouterGroup
	UserService
}

func UseUserHandler(params NewUserHandlerParams) {
	handler := &UserHandler{
		userService: params.UserService,
	}

	router := params.Router.Group("/users")

	router.GET("/:id", handler.GetUserById)
}

func (h *UserHandler) GetUserById(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	resp, err := h.userService.GetUserById(int64(userId))
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, http.StatusOK, resp)
}

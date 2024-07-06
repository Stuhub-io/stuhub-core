package rest

import (
	"net/http"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/services/user"
	"github.com/Stuhub-io/internal/api/decorators"
	"github.com/Stuhub-io/internal/api/middleware"
	"github.com/Stuhub-io/internal/api/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *user.Service
}

type NewUserHandlerParams struct {
	Router         *gin.RouterGroup
	AuthMiddleware *middleware.AuthMiddleware
	UserService    *user.Service
}

func UseUserHandler(params NewUserHandlerParams) {
	handler := &UserHandler{
		userService: params.UserService,
	}

	router := params.Router.Group("/user-services")
	authMiddleware := params.AuthMiddleware

	router.Use(authMiddleware.Authenticated())

	router.GET("/:id", decorators.CurrentUser(handler.GetUserById))
	router.GET("/email/:email", handler.GetUserByEmail)
}

// GetUserByID godoc
//
//	@Summary		Get User Details
//	@Description	Get User Details by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	domain.User
//	@Failure		400	{object}	domain.Error
//	@Failure		500	{object}	domain.Error
//	@Router			/v1/user-services/{id} [get]
func (h *UserHandler) GetUserById(c *gin.Context, user *domain.User) {
	// userId := c.Param("id")
	// _, err := h.userService.GetUserById(userId)
	// if err != nil {
	// 	response.WithErrorMessage(c, err.Code, err.Error, err.Message)
	// 	return
	// }

	response.WithData(c, http.StatusOK, user)
}

func (h *UserHandler) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	resp, err := h.userService.GetUserByEmail(email)
	if err != nil {
		response.WithErrorMessage(c, err.Code, err.Error, err.Message)
		return
	}

	response.WithData(c, http.StatusOK, resp)
}

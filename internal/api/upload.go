package api

import (
	"net/http"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/services/upload"
	"github.com/Stuhub-io/internal/api/decorators"
	"github.com/Stuhub-io/internal/api/middleware"
	"github.com/Stuhub-io/internal/api/request"
	"github.com/Stuhub-io/internal/api/response"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	uploadService *upload.UploadService
}

type NewUploadHandlerParams struct {
	Router         *gin.RouterGroup
	AuthMiddleware *middleware.AuthMiddleware
	UploadService  *upload.UploadService
}

func UseUploadHandler(params NewUploadHandlerParams) {
	handler := &UploadHandler{
		uploadService: params.UploadService,
	}

	router := params.Router.Group("/upload-services")
	authMiddleware := params.AuthMiddleware

	router.Use(authMiddleware.Authenticated())

	router.POST("/sign-url", decorators.CurrentUser(handler.SignUploadUrl))
}

func (h *UploadHandler) SignUploadUrl(c *gin.Context, user *domain.User) {
	var body request.SignUrlRequestBody
	if err := request.Validate(c, &body); err != nil {
		response.BindError(c, err.Error())
		return
	}
	signedData, err := h.uploadService.GenerateSignedUrl(domain.SignUrlInput{
		PublicID:        body.PublicID,
		ResourceType:    body.ResourceType,
		AdditionalQuery: body.AdditionalQuery,
	})
	if err != nil {
		response.BindError(c, err.Message)
		return
	}
	response.WithData(c, http.StatusOK, signedData)
}

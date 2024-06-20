package response

import (
	"net/http"

	"github.com/Stuhub-io/core/domain"
	"github.com/gin-gonic/gin"
)

type MessageResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type DataResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data"`
}

type PaginationResponse struct {
	Code        int    `json:"code"`
	Message     string `json:"message,omitempty"`
	Data        any    `json:"data"`
	HasNextPage bool   `json:"has_next"`
	Count       int    `json:"count"`
}

type PaginationPayload struct {
	Data        any  `json:"data"`
	HasNextPage bool `json:"has_next"`
	Count       int  `json:"count"`
}

func WithMessage(c *gin.Context, code int, message string) {
	c.JSON(code, &MessageResponse{
		Code:    code,
		Message: message,
	})
}

func WithData(c *gin.Context, code int, data any, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(code, &DataResponse{
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

func WithPagination(c *gin.Context, code int, data PaginationPayload, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}

	c.JSON(code, &PaginationResponse{
		Code:        code,
		Message:     msg,
		Data:        data.Data,
		HasNextPage: data.HasNextPage,
		Count:       data.Count,
	})
}

func WithErrorMessage(c *gin.Context, code int, err string, message string) {
	c.JSON(code, &ErrorResponse{
		Code:    code,
		Error:   err,
		Message: message,
	})
}

func ServerError(c *gin.Context) {
	message := "the server encountered a problem and could not process your request"

	WithErrorMessage(c, http.StatusInternalServerError, domain.ErrInternalServerError.Error(), message)
}

func NotFound(c *gin.Context) {
	message := "the requested resource could not be found"

	WithErrorMessage(c, http.StatusNotFound, domain.ErrNotFound.Error(), message)
}

func BadRequest(c *gin.Context, err error, message ...string) {
	msg := "the server cannot understand or process correctly"
	if len(message) > 0 {
		msg = message[0]
	}

	WithErrorMessage(c, http.StatusBadRequest, domain.ErrBadRequest.Error(), msg)
}

func Unauthorized(c *gin.Context) {
	message := "you are not authorized"

	WithErrorMessage(c, http.StatusUnauthorized, domain.ErrUnauthorized.Error(), message)
}

func BindError(c *gin.Context, error string) {
	message := "please provide the correct body input"

	WithErrorMessage(c, http.StatusBadRequest, domain.ErrBadParamInput.Error(), message)
}

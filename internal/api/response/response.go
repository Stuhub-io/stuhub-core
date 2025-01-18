package response

import (
	"github.com/Stuhub-io/core/domain"
	"github.com/gin-gonic/gin"
)

type MessageResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

type DataResponse struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data"`
}

type PaginationResponse[T any] struct {
	Status     string `json:"status"`
	Code       int    `json:"code"`
	Message    string `json:"message,omitempty"`
	Data       any    `json:"data"`
	Pagination T      `json:"pagination"`
}

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

func WithMessage(c *gin.Context, code int, message string) {
	c.JSON(code, &MessageResponse{
		Status:  StatusSuccess,
		Code:    code,
		Message: message,
	})
}

func WithData(c *gin.Context, code int, data any, message ...string) {
	msg := getMessage("", message...)

	c.JSON(code, &DataResponse{
		Status:  StatusSuccess,
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

func WithPagination(c *gin.Context, code int, data any, pagination domain.Pagination, message ...string) {
	msg := getMessage("", message...)

	c.JSON(code, &PaginationResponse[domain.Pagination]{
		Status:     StatusSuccess,
		Code:       code,
		Message:    msg,
		Data:       data,
		Pagination: pagination,
	})
}

func WithCursorPagination[K domain.Cursor](c *gin.Context, code int, data any, pagination domain.CursorPagination[K], message ...string) {
	msg := getMessage("", message...)

	c.JSON(code, &PaginationResponse[domain.CursorPagination[K]]{
		Status:     StatusSuccess,
		Code:       code,
		Message:    msg,
		Data:       data,
		Pagination: pagination,
	})
}

func WithErrorMessage(c *gin.Context, code int, err string, message string) {
	c.JSON(code, &ErrorResponse{
		Status:  StatusError,
		Code:    code,
		Error:   err,
		Message: message,
	})
}

func OK(c *gin.Context, message ...string) {
	msg := getMessage(domain.SuccessOK.Message, message...)

	WithMessage(c, domain.SuccessOK.Code, msg)
}

func Created(c *gin.Context, message ...string) {
	msg := getMessage(domain.SuccessCreated.Message, message...)

	WithMessage(c, domain.SuccessCreated.Code, msg)
}

func ServerError(c *gin.Context, message ...string) {
	msg := getMessage(domain.ErrInternalServerError.Message, message...)

	WithErrorMessage(c, domain.ErrInternalServerError.Code, domain.ErrInternalServerError.Error, msg)
}

func Forbidden(c *gin.Context, message ...string) {
	msg := getMessage(domain.ErrForbidden.Message, message...)

	WithErrorMessage(c, domain.ErrForbidden.Code, domain.ErrForbidden.Error, msg)
}

func NotFound(c *gin.Context, message ...string) {
	msg := getMessage(domain.ErrNotFound.Message, message...)

	WithErrorMessage(c, domain.ErrNotFound.Code, domain.ErrNotFound.Error, msg)
}

func BadRequest(c *gin.Context, message ...string) {
	msg := getMessage(domain.ErrBadRequest.Message, message...)

	WithErrorMessage(c, domain.ErrBadRequest.Code, domain.ErrBadRequest.Error, msg)
}

func Unauthorized(c *gin.Context, message ...string) {
	msg := getMessage(domain.ErrUnauthorized.Message, message...)

	WithErrorMessage(c, domain.ErrUnauthorized.Code, domain.ErrUnauthorized.Error, msg)
}

func Conflict(c *gin.Context, message ...string) {
	msg := getMessage(domain.ErrConflict.Message, message...)

	WithErrorMessage(c, domain.ErrConflict.Code, domain.ErrConflict.Error, msg)
}

func BindError(c *gin.Context, err string) {
	WithErrorMessage(c, domain.ErrBadParamInput.Code, domain.ErrBadParamInput.Error, err)
}

func getMessage(base string, message ...string) string {
	msg := base
	if len(message) > 0 {
		msg = message[0]
	}

	return msg
}

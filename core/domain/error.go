package domain

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code    int    `json:"code"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

const (
	InternalServerErrCode = http.StatusInternalServerError
	NotFoundCode          = http.StatusNotFound
	BadRequestCode        = http.StatusBadRequest
	UnauthorizedCode      = http.StatusUnauthorized
	ConflictCode          = http.StatusConflict
	ForbiddenCode         = http.StatusForbidden
)

const (
	InternalServerErr = "Unexpected server error."
	DatabaseErr       = "Unexpected database error."
	NotFoundErr       = "The requested resource was not found"
	BadRequestErr     = "Bad request."
	UnauthorizedErr   = "Unauthorized access."
	ConflictErr       = "Conflict occurred."
	ForbiddenErr      = "Access forbidden."
)

var (
	ErrInternalServerError = &Error{
		Code:    InternalServerErrCode,
		Error:   InternalServerErr,
		Message: "The server encountered a problem and could not process your request",
	}
	ErrNotFound = &Error{
		Code:    NotFoundCode,
		Error:   NotFoundErr,
		Message: "The requested resource could not be found",
	}
	ErrBadRequest = &Error{
		Code:    BadRequestCode,
		Error:   BadRequestErr,
		Message: "You are not authorized",
	}
	ErrUnauthorized = &Error{
		Code:    UnauthorizedCode,
		Error:   UnauthorizedErr,
		Message: "The server cannot understand or process correctly",
	}
	ErrConflict = &Error{
		Code:    ConflictCode,
		Error:   ConflictErr,
		Message: "The request could not be completed due to a conflict with the current state of the resource",
	}
	ErrForbidden = &Error{
		Code:    ForbiddenCode,
		Error:   ForbiddenErr,
		Message: "You do not have permission to access this resource",
	}
	ErrBadParamInput = &Error{
		Code:    BadRequestCode,
		Error:   BadRequestErr,
		Message: "The request contains invalid parameters. Please check your input and try again",
	}
)

var (
	ErrDatabaseQuery = &Error{
		Code:    InternalServerErrCode,
		Error:   DatabaseErr,
		Message: "Database can't not process the query",
	}
	ErrDatabaseMutation = &Error{
		Code:    InternalServerErrCode,
		Error:   DatabaseErr,
		Message: "Database can't not process the mutation",
	}
)

var (
	ErrUserNotFoundById = func(id int) *Error {
		return &Error{
			Code:    NotFoundCode,
			Error:   NotFoundErr,
			Message: fmt.Sprintf("The user with the ID '%d' does not exist.", id),
		}
	}
	ErrUserNotFoundByEmail = func(email string) *Error {
		return &Error{
			Code:    NotFoundCode,
			Error:   NotFoundErr,
			Message: fmt.Sprintf("The user with the Email '%s' does not exist.", email),
		}
	}
	ErrExistUserEmail = func(email string) *Error {
		return &Error{
			Code:    BadRequestCode,
			Error:   BadRequestErr,
			Message: fmt.Sprintf("The user with the Email '%s' already exist.", email),
		}
	}
)

var (
	ErrInvalidCredentials = &Error{
		Code:    BadRequestCode,
		Error:   BadRequestErr,
		Message: "Invalid credentials. Please input the correct account!",
	}
)

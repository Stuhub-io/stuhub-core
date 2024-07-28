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
		Message: "The server cannot understand or process correctly",
	}
	ErrUnauthorized = &Error{
		Code:    UnauthorizedCode,
		Error:   UnauthorizedErr,
		Message: "You are not authorized",
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
		Message: "Database can't process the query",
	}
	ErrDatabaseMutation = &Error{
		Code:    InternalServerErrCode,
		Error:   DatabaseErr,
		Message: "Database can't process the mutation",
	}
	ErrRollbackErr = &Error{
		Code:    InternalServerErrCode,
		Error:   DatabaseErr,
		Message: "Database can't rollback the transaction",
	}
)

var (
	ErrUserPassword = &Error{
		Code:    BadRequestCode,
		Error:   BadRequestErr,
		Message: "The password is invalid. Please input the correct password!",
	}
	ErrUserNotFound = &Error{
		Code:    NotFoundCode,
		Error:   NotFoundErr,
		Message: "The user does not exist.",
	}
	ErrUserNotFoundById = func(id string) *Error {
		return &Error{
			Code:    NotFoundCode,
			Error:   NotFoundErr,
			Message: fmt.Sprintf("The user with the ID '%s' does not exist.", id),
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

var (
	ErrTokenExpired = &Error{
		Code:    BadRequestCode,
		Error:   BadRequestErr,
		Message: "The resource was expired. Please request a new one!",
	}
)

var (
	ErrSendMail = &Error{
		Code:    InternalServerErrCode,
		Error:   InternalServerErr,
		Message: "Failed to process sending email. Please try again!",
	}
)

var (
	ErrGetGoogleInfo = &Error{
		Code:    InternalServerErrCode,
		Error:   InternalServerErr,
		Message: "Failed to get Google info. Please try again!",
	}
)

var (
	ErrOrgNotFound = &Error{
		Code:    NotFoundCode,
		Error:   NotFoundErr,
		Message: "The organization does not exist.",
	}
	ErrExistOwnerOrg = func(name string) *Error {
		return &Error{
			Code:    BadRequestCode,
			Error:   BadRequestErr,
			Message: fmt.Sprintf("The org with the Name '%s' already exist.", name),
		}
	}
	ErrOrgMemberNotFound = &Error{
		Code:    NotFoundCode,
		Error:   NotFoundErr,
		Message: "The member does not exist.",
	}
	ErrExistOrgMember = func(userPkID int64) *Error {
		return &Error{
			Code:    BadRequestCode,
			Error:   BadRequestErr,
			Message: fmt.Sprintf("The member with the ID '%d' already exist in this organization.", userPkID),
		}
	}
)

func NewErr(msg string, code int) *Error {
	return &Error{
		Code:    code,
		Error:   http.StatusText(code),
		Message: msg,
	}
}

package domain

import (
	"net/http"
)

type Success struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	OKCode      = http.StatusOK
	CreatedCode = http.StatusAccepted
)

var (
	SuccessOK = &Success{
		Code:    OKCode,
		Message: "Request proceeded successfully",
	}
	SuccessCreated = &Success{
		Code:    CreatedCode,
		Message: "Resource successfully created",
	}
)

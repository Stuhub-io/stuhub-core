package rest

import (
	"errors"
	"net/http"

	"github.com/Stuhub-io/core/domain"
	"github.com/sirupsen/logrus"
)

type ResponseError struct {
	Message string `json:"message"`
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	if errors.Is(err, domain.ErrInternalServerError) {
		return http.StatusInternalServerError
	}
	if errors.Is(err, domain.ErrNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, domain.ErrConflict) {
		return http.StatusConflict
	}

	return http.StatusInternalServerError
}

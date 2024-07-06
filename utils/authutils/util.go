package authutils

import (
	"errors"
	"strings"
)

type userPayload string

var (
	UserPayload userPayload = "userPayload"
)

func ExtractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	token := strings.Split(header, " ")
	if len(token) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return token[1], nil
}

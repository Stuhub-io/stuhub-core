package authutils

import (
	"errors"
	"strings"
)

type tokenPayloadKey string

var (
	TokenPayloadKey tokenPayloadKey = "payload"
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

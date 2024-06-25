package token

import (
	"fmt"
	"time"

	"github.com/Stuhub-io/core/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 5

type JWTMaker struct {
	secretKey string
}

func new(secretKey string) (*JWTMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

func Must(secretKey string) *JWTMaker {
	jwtMaker, err := new(secretKey)
	if err != nil {
		panic(err)
	}

	return jwtMaker
}

func (m *JWTMaker) CreateToken(email string, duration time.Duration) (string, error) {
	claims, err := newPayload(email, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(m.secretKey))
}

func (m *JWTMaker) VerifyToken(token string) (*domain.TokenPayload, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token: ")
		}

		return []byte(m.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token: ")
	}

	return &domain.TokenPayload{
		ID:        claims.Subject,
		Email:     claims.Issuer,
		IssuedAt:  claims.IssuedAt.Local(),
		ExpiredAt: claims.ExpiresAt.Local(),
	}, nil
}

func newPayload(email string, duration time.Duration) (*jwt.RegisteredClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	claims := &jwt.RegisteredClaims{
		Subject:   tokenID.String(),
		Issuer:    email,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
	}

	return claims, nil
}

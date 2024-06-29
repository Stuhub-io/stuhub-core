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
type CustomClaims struct {
	jwt.RegisteredClaims
	UserPkID int64  `json:"user_pkid,string"`
	Email    string `json:"email"`
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

func (m *JWTMaker) CreateToken(pkid int64, email string, duration time.Duration) (string, error) {
	claims, err := newPayload(pkid, email, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(m.secretKey))
}

func (m *JWTMaker) DecodeToken(token string) (*domain.TokenPayload, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token: ")
		}

		return []byte(m.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &CustomClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*CustomClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token: ")
	}

	return &domain.TokenPayload{
		UserPkID:  claims.UserPkID,
		Email:     claims.Email,
		IssuedAt:  claims.RegisteredClaims.IssuedAt.Local(),
		ExpiredAt: claims.RegisteredClaims.ExpiresAt.Local(),
	}, nil
}

func newPayload(pkid int64, email string, duration time.Duration) (*CustomClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	claims := &CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   tokenID.String(),
			Issuer:    email,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		UserPkID: pkid,
		Email:    email,
	}
	return claims, nil
}

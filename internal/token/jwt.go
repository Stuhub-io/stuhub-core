package token

import (
	"errors"
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

type CustomAuthClaims struct {
	jwt.RegisteredClaims
	UserPkID int64  `json:"user_pkid,string"`
	Email    string `json:"email"`
}

type CustomOrgInviteClaims struct {
	jwt.RegisteredClaims
	UserPkID int64 `json:"user_pkid"`
	OrgPkID  int64 `json:"org_pkid"`
}

func newMaker(secretKey string) (*JWTMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	return &JWTMaker{secretKey}, nil
}

func Must(secretKey string) *JWTMaker {
	jwtMaker, err := newMaker(secretKey)
	if err != nil {
		panic(err)
	}

	return jwtMaker
}

func (m *JWTMaker) CreateToken(pkid int64, email string, duration time.Duration) (string, error) {
	claims, err := newAuthPayload(pkid, email, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return jwtToken.SignedString([]byte(m.secretKey))
}

func (m *JWTMaker) DecodeToken(token string) (*domain.TokenAuthPayload, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(m.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &CustomAuthClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*CustomAuthClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return &domain.TokenAuthPayload{
		UserPkID:  claims.UserPkID,
		Email:     claims.Email,
		IssuedAt:  claims.RegisteredClaims.IssuedAt.Local(),
		ExpiredAt: claims.RegisteredClaims.ExpiresAt.Local(),
	}, nil
}

func (m *JWTMaker) CreateOrgInviteToken(userPkID, orgPkID int64, duration time.Duration) (string, error) {
	claims, err := newOrgInvitePayload(userPkID, orgPkID, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return jwtToken.SignedString([]byte(m.secretKey))
}

func (m *JWTMaker) DecodeOrgInviteToken(token string) (*domain.TokenOrgInvitePayload, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(m.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &CustomOrgInviteClaims{}, keyFunc)
	if err != nil {
		return nil, err
	}

	claims, ok := jwtToken.Claims.(*CustomOrgInviteClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return &domain.TokenOrgInvitePayload{
		UserPkID:  claims.UserPkID,
		OrgPkID:   claims.OrgPkID,
		IssuedAt:  claims.RegisteredClaims.IssuedAt.Local(),
		ExpiredAt: claims.RegisteredClaims.ExpiresAt.Local(),
	}, nil
}

func newAuthPayload(pkid int64, email string, duration time.Duration) (*CustomAuthClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	claims := &CustomAuthClaims{
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

func newOrgInvitePayload(userPkID, orgPkID int64, duration time.Duration) (*CustomOrgInviteClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	claims := &CustomOrgInviteClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   tokenID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
		UserPkID: userPkID,
		OrgPkID:  orgPkID,
	}

	return claims, nil
}

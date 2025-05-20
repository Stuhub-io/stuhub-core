package middleware

import (
	"context"

	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/utils/authutils"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tokenMaker ports.TokenMaker
	repo       *ports.Repository
}

type NewAuthMiddlewareParams struct {
	ports.TokenMaker
	*ports.Repository
}

func NewAuthMiddleware(params NewAuthMiddlewareParams) *AuthMiddleware {
	return &AuthMiddleware{
		tokenMaker: params.TokenMaker,
		repo:       params.Repository,
	}
}

func (a *AuthMiddleware) Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := authutils.ExtractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.Next()
			return
		}

		payload, err := a.tokenMaker.DecodeToken(token)
		if err != nil {
			c.Next()
			return
		}

		user, dbErr := a.repo.User.GetUserByPkID(context.Background(), payload.UserPkID)

		if dbErr != nil {
			c.Next()
			return
		}

		c.Set(string(authutils.UserPayloadKey), user)

		c.Next()
	}
}

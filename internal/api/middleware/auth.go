package middleware

import (
	"context"

	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/utils/authutils"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tokenMaker     ports.TokenMaker
	userRepository ports.UserRepository
}

type NewAuthMiddlewareParams struct {
	ports.TokenMaker
	ports.UserRepository
}

func NewAuthMiddleware(params NewAuthMiddlewareParams) *AuthMiddleware {
	return &AuthMiddleware{
		tokenMaker:     params.TokenMaker,
		userRepository: params.UserRepository,
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

		user, dbErr := a.userRepository.GetUserByPkID(context.Background(), payload.UserPkID)

		if dbErr != nil {
			c.Next()
			return
		}

		c.Set(string(authutils.UserPayloadKey), user)

		c.Next()
	}
}

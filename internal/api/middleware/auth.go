package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/utils/authutils"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	tokenMaker     ports.TokenMaker
	userRepository ports.UserRepository
	cacheStore     ports.CacheStore
}

type NewAuthMiddlewareParams struct {
	ports.TokenMaker
	ports.UserRepository
	ports.CacheStore
}

func NewAuthMiddleware(params NewAuthMiddlewareParams) *AuthMiddleware {
	return &AuthMiddleware{
		tokenMaker:     params.TokenMaker,
		userRepository: params.UserRepository,
		cacheStore:     params.CacheStore,
	}
}

func (a *AuthMiddleware) Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := authutils.ExtractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(domain.UnauthorizedCode, domain.ErrUnauthorized)
			return
		}

		payload, err := a.tokenMaker.DecodeToken(token)
		if err != nil {
			c.AbortWithStatusJSON(domain.UnauthorizedCode, domain.ErrUnauthorized)
			return
		}

		var user *domain.User
		userPkID := payload.UserPkID
		user = a.cacheStore.GetUser(userPkID)
		if user == nil {
			data, err := a.userRepository.GetUserByPkID(context.Background(), userPkID)
			if err != nil {
				c.AbortWithStatusJSON(err.Code, err)
				return
			}

			cacheError := a.cacheStore.SetUser(data, time.Hour)
			if cacheError != nil {
				fmt.Printf("caching error: %v", cacheError)
			}

			user = data
		}

		c.Set(string(authutils.UserPayloadKey), user)

		c.Next()
	}
}

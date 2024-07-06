package middleware

import (
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/Stuhub-io/utils/authutils"
	"github.com/gin-gonic/gin"
)

type HandlerWithTokenPayload func(*gin.Context, ...*domain.TokenPayload)

func Authenticated(tokenMaker ports.TokenMaker) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := authutils.ExtractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(domain.UnauthorizedCode, domain.ErrUnauthorized)
			return
		}

		payload, err := tokenMaker.DecodeToken(token)
		if err != nil {
			c.AbortWithStatusJSON(domain.UnauthorizedCode, domain.ErrUnauthorized)
			return
		}

		c.Set(string(authutils.TokenPayloadKey), payload)

		c.Next()
	}
}

package decorators

import (
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/utils/authutils"
	"github.com/gin-gonic/gin"
)

type HandlerWithCurrentUser func(*gin.Context, *domain.User)

func CurrentUser(handler HandlerWithCurrentUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, Ok := c.Keys[string(authutils.UserPayloadKey)].(*domain.User)

		if !Ok {
			handler(c, nil)
			return
		}

		handler(c, currentUser)
	}
}

func RequiredAuth(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, Ok := c.Keys[string(authutils.UserPayloadKey)].(*domain.User); !Ok {
			c.AbortWithStatusJSON(domain.UnauthorizedCode, domain.ErrUnauthorized)
			return
		}

		handler(c)
	}
}

package decorators

import (
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/utils/authutils"
	"github.com/gin-gonic/gin"
)

type HandlerWithCurrentUser func(*gin.Context, *domain.User)

func CurrentUser(handler HandlerWithCurrentUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUser, Ok := c.Keys[string(authutils.UserPayload)].(*domain.User)
		if !Ok {
			c.AbortWithStatusJSON(domain.UnauthorizedCode, domain.ErrUnauthorized)
			return
		}

		handler(c, currentUser)
	}
}

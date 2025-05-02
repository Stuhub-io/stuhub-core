package middleware

import (
	"fmt"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

type ServiceAuthMiddleware struct {
	ServiceKeys []string
}

func NewServiceAuthMiddleware(serviceKeys []string) *ServiceAuthMiddleware {
	return &ServiceAuthMiddleware{
		ServiceKeys: serviceKeys,
	}
}

func (m *ServiceAuthMiddleware) RequiredServiceKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		fmt.Print("\n\nService key: ", token, "\n\n")

		authenticated := slices.Contains(m.ServiceKeys, token)

		if !authenticated {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid service credentials"})
			return
		}

		// Add service info to the context
		c.Set("is_internal_service", true)

		c.Next()
	}
}

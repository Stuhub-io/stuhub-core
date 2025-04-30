package middleware

import (
	"fmt"
	"net/http"

	"github.com/Stuhub-io/utils/authutils"
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
		token, err := authutils.ExtractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "service authentication required"})
			return
		}

		fmt.Print("Service key: ", token, "\n")

		var authenticated bool

		for _, apiKey := range m.ServiceKeys {
			if token == apiKey {
				authenticated = true
				break
			}
		}

		if !authenticated {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid service credentials"})
			return
		}

		// Add service info to the context
		c.Set("is_internal_service", true)

		c.Next()
	}
}

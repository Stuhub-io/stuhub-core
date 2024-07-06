package middleware

import (
	"github.com/Stuhub-io/config"
	"github.com/gin-gonic/gin"
)

func JSON(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Next()
	}
}

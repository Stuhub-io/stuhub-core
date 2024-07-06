package middleware

import (
	"github.com/Stuhub-io/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(cfg *config.Config) gin.HandlerFunc {

	return func(c *gin.Context) {
		cors.New(
			cors.Config{
				AllowOrigins: cfg.AllowedOrigins,
				AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
				AllowHeaders: []string{
					"Origin", "Host", "Content-Type", "Content-Length", "Accept-Encoding", "Accept-Language", "Accept",
					"X-CSRF-Token", "Authorization", "X-Requested-With", "X-Access-Token", "credentials",
				},
				AllowCredentials: true,
			},
		)(c)
	}
}

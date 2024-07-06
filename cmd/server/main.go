package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/services/auth"
	"github.com/Stuhub-io/core/services/user"
	_ "github.com/Stuhub-io/docs"
	api "github.com/Stuhub-io/internal/api"
	"github.com/Stuhub-io/internal/api/middleware"
	"github.com/Stuhub-io/internal/cache"
	"github.com/Stuhub-io/internal/cache/redis"
	"github.com/Stuhub-io/internal/hasher"
	"github.com/Stuhub-io/internal/mailer"
	"github.com/Stuhub-io/internal/remote"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/postgres"
	"github.com/Stuhub-io/internal/token"
	"github.com/Stuhub-io/logger"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:5000
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	cfg := config.LoadConfig(config.GetDefaultConfigLoaders())

	logger := logger.NewLogrusLogger()

	postgresDB := postgres.Must(cfg.DBDsn, cfg.Debug)

	redisCache := redis.Must(cfg.CacheHost, cfg.CachePort, cfg.CachePassword)

	tokenMaker := token.Must(cfg.SecretKey)

	dbStore := store.NewDBStore(postgresDB, redisCache)
	cacheStore := cache.NewCacheStore(redisCache)

	hasher := hasher.NewScrypt([]byte(cfg.HashPwSecretKey))

	// TODO: read from env
	mailer := mailer.NewMailer(mailer.NewMailerParams{
		Name:      "",
		Address:   "",
		ClientKey: "",
		Config:    cfg,
	})

	r := gin.Default()

	r.Use(middleware.CORS(&cfg))
	r.Use(middleware.JSON(&cfg))

	remoteRoute := remote.NewRemoteRoute()

	// repositories
	userRepository := postgres.NewUserRepository(postgres.NewUserRepositoryParams{
		Store: dbStore,
		Cfg:   cfg,
	})

	// services
	userService := user.NewService(user.NewServiceParams{
		UserRepository: userRepository,
		Config:         cfg,
	})
	authService := auth.NewService(auth.NewServiceParams{
		UserRepository: userRepository,
		TokenMaker:     tokenMaker,
		Mailer:         mailer,
		Config:         cfg,
		RemoteRoute:    remoteRoute,
		Hasher:         hasher,
	})

	authMiddleware := middleware.NewAuthMiddleware(middleware.NewAuthMiddlewareParams{
		TokenMaker:     tokenMaker,
		UserRepository: userRepository,
		CacheStore:     cacheStore,
	})

	// handlers
	v1 := r.Group("/v1")
	{
		api.UseUserHandler(api.NewUserHandlerParams{
			Router:         v1,
			AuthMiddleware: authMiddleware,
			UserService:    userService,
		})
		api.UseAuthHandler(api.NewAuthHandlerParams{
			Router:      v1,
			AuthService: authService,
		})
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "Server is running OK!")
	})

	r.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: r,
	}

	go func() {
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Server closed!")
		} else if err != nil {
			fmt.Printf("Error starting server at %s\n", err)
			os.Exit(1)
		}
	}()

	logger.Infof("Server started at %s", fmt.Sprintf("port %d", cfg.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit

	shutdownServer(srv, cfg.GetShutdownTimeout(), logger)
}

func shutdownServer(srv *http.Server, timeout time.Duration, logger logger.Logger) {
	logger.Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(err, "failed to shutdown server")
	}

	logger.Info("Server exiting")
}

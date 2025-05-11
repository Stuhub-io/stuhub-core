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
	"github.com/Stuhub-io/core/services/activity"
	"github.com/Stuhub-io/core/services/activity_v2"
	"github.com/Stuhub-io/core/services/auth"
	"github.com/Stuhub-io/core/services/organization"
	"github.com/Stuhub-io/core/services/page"
	pageAccessLog "github.com/Stuhub-io/core/services/page_access_log"
	"github.com/Stuhub-io/core/services/upload"
	"github.com/Stuhub-io/core/services/user"
	_ "github.com/Stuhub-io/docs"
	"github.com/Stuhub-io/internal/api"
	"github.com/Stuhub-io/internal/api/middleware"
	"github.com/Stuhub-io/internal/cache"
	"github.com/Stuhub-io/internal/cache/redis"
	"github.com/Stuhub-io/internal/hasher"
	"github.com/Stuhub-io/internal/mailer"
	"github.com/Stuhub-io/internal/oauth"
	"github.com/Stuhub-io/internal/remote"
	store "github.com/Stuhub-io/internal/repository"
	"github.com/Stuhub-io/internal/repository/postgres"
	"github.com/Stuhub-io/internal/repository/scylla"
	"github.com/Stuhub-io/internal/token"
	"github.com/Stuhub-io/internal/uploader"
	"github.com/Stuhub-io/logger"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

type TempCache struct{}

func (TempCache) Set(key string, value any, duration time.Duration) error { return nil }
func (TempCache) Get(key string) (string, error)                          { return "", nil }
func (TempCache) Delete(key string) error {
	return nil
}

func main() {
	cfg := config.LoadConfig(config.GetDefaultConfigLoaders())

	logger := logger.NewLogrusLogger()

	postgresDB := postgres.Must(cfg.DBDsn, cfg.Debug, logger)

	scyllaDB := scylla.Must(cfg.ScyllaHosts, cfg.ScyllaPort, cfg.ScyllaKeyspace, cfg.Debug, logger)

	// redisCache := redis.Must(cfg.RedisUrl, logger)
	redisCache := redis.Mock()

	// elasticSearch := elasticsearch.Must(cfg.ElasticSearchURL, logger)

	tokenMaker := token.Must(cfg.SecretKey)

	cacheStore := cache.NewCacheStore(redisCache)
	dbStore := store.NewDBStore(postgresDB, cacheStore, scyllaDB)

	hasher := hasher.NewScrypt([]byte(cfg.HashPwSecretKey))

	mailer := mailer.NewMailer(mailer.NewMailerParams{
		From:      "Stuhub.IO",
		Address:   cfg.SendgridEmailFrom,
		ClientKey: cfg.SendgridKey,
		Logger:    logger,
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
	orgRepository := postgres.NewOrganizationRepository(postgres.NewOrganizationRepositoryParams{
		Store:          dbStore,
		Cfg:            cfg,
		UserRepository: userRepository,
	})
	pageRepository := postgres.NewPageRepository(postgres.NewPageRepositoryParams{
		Cfg:   cfg,
		Store: dbStore,
	})
	organizationInviteRepository := postgres.NewOrganizationInvitesRepository(
		postgres.NewOrganizationInvitesRepositoryParams{
			Cfg:   cfg,
			Store: dbStore,
		},
	)
	pageAccessLogsRepository := postgres.NewPageAccessLogRepository(
		postgres.NewPageAccessLogRepositoryParams{
			Cfg:   cfg,
			Store: dbStore,
		},
	)

	activityV2Repository := postgres.NewActivityV2Repository(
		postgres.ActivityV2RepositoryParams{
			Cfg:   cfg,
			Store: dbStore,
		},
	)

	activityRepository := scylla.NewActivityRepository(scylla.ActivityRepositoryParams{
		Cfg:   cfg,
		Store: dbStore,
	})

	// indexers
	// pageIndexer := elasticsearch.NewPageIndexer(elasticSearch)
	// pageIndexer := elasticsearch.NewPageIndexer(nil)
	// pageIndexer.Index(context.Background())

	// services
	cloudinaryUploader := uploader.NewCloudinaryUploader(cfg)
	authMiddleware := middleware.NewAuthMiddleware(middleware.NewAuthMiddlewareParams{
		TokenMaker:     tokenMaker,
		UserRepository: userRepository,
	})

	serviceAuthMiddleware := middleware.NewServiceAuthMiddleware(cfg.InternalServiceApiKeys)

	oauthService := oauth.NewOauthService(logger)
	userService := user.NewService(user.NewServiceParams{
		Config:         cfg,
		UserRepository: userRepository,
	})
	authService := auth.NewService(auth.NewServiceParams{
		Config:         cfg,
		UserRepository: userRepository,
		PageRepository: pageRepository,
		OauthService:   oauthService,
		TokenMaker:     tokenMaker,
		Mailer:         mailer,
		RemoteRoute:    remoteRoute,
		Hasher:         hasher,
	})
	orgService := organization.NewService(organization.NewServiceParams{
		Config:                       cfg,
		OrganizationRepository:       orgRepository,
		UserRepository:               userRepository,
		TokenMaker:                   tokenMaker,
		Hasher:                       hasher,
		Mailer:                       mailer,
		RemoteRoute:                  remoteRoute,
		OrganizationInviteRepository: organizationInviteRepository,
	})
	pageService := page.NewService(page.NewServiceParams{
		Config:                  cfg,
		Logger:                  logger,
		PageRepository:          pageRepository,
		PageAccessLogRepository: pageAccessLogsRepository,
		Mailer:                  mailer,
		ActivityRepository:      activityRepository,
		ActivityV2Repository:    activityV2Repository,
	})
	uploadService := upload.NewUploadService(upload.NewUploadServiceParams{
		Config:   cfg,
		Uploader: cloudinaryUploader,
	})
	pageAccessLogService := pageAccessLog.NewService(pageAccessLog.NewServiceParams{
		PageRepository:          pageRepository,
		PageAccessLogRepository: pageAccessLogsRepository,
	})

	activityService := activity.NewService(activity.NewServiceParams{
		PageRepository:         pageRepository,
		Config:                 cfg,
		Logger:                 logger,
		OrganizationRepository: orgRepository,
		ActivityRepository:     activityRepository,
		UserRepository:         userRepository,
	})

	activityV2Service := activity_v2.NewService(activity_v2.NewServiceParams{
		Config:                 cfg,
		Logger:                 logger,
		PageRepository:         pageRepository,
		OrganizationRepository: orgRepository,
		ActivityV2Repository:   activityV2Repository,
		UserRepository:         userRepository,
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
		api.UseOrganizationHandler(api.NewOrganizationHandlerParams{
			Router:         v1,
			AuthMiddleware: authMiddleware,
			OrgService:     orgService,
		})
		api.UsePageHandle((api.NewPageHandlerParams{
			Router:                v1,
			AuthMiddleware:        authMiddleware,
			PageService:           pageService,
			ServiceAuthMiddleware: serviceAuthMiddleware,
		}))
		api.UseUploadHandler(api.NewUploadHandlerParams{
			Router:         v1,
			AuthMiddleware: authMiddleware,
			UploadService:  uploadService,
		})
		api.UsePageAccessLogHandler(api.NewPageAccessLogHandlerParams{
			Router:               v1,
			AuthMiddleware:       authMiddleware,
			PageAccessLogService: pageAccessLogService,
		})
		api.UseActivityHandler(api.NewActivityHandlerParams{
			Router:          v1,
			AuthMiddleware:  authMiddleware,
			ActivityService: activityService,
		})
	}

	v2 := r.Group("/v2")
	{
		api.UseActivityV2Handler(api.NewActivityV2HandlerParams{
			Router:            v2,
			AuthMiddleware:    authMiddleware,
			ActivityV2Service: activityV2Service,
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

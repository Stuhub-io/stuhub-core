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
	"github.com/Stuhub-io/internal/mailer"
	"github.com/Stuhub-io/internal/repository/postgres"
	"github.com/Stuhub-io/internal/rest"
	"github.com/Stuhub-io/internal/rest/middleware"
	"github.com/Stuhub-io/internal/token"
	"github.com/Stuhub-io/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig(config.GetDefaultConfigLoaders())

	logger := logger.NewLogrusLogger()

	postgresDB, err := postgres.NewStore(cfg.DBDsn, cfg.Debug)
	if err != nil {
		panic(err)
	}

	tokenMaker := token.Must(cfg.SecretKey)

	// TODO: read from env
	mailer := mailer.NewMailer(mailer.NewMailerParams{
		Name:      "",
		Address:   "",
		ClientKey: "",
		Config:    cfg,
	})

	r := gin.Default()

	r.Use(middleware.CORS())

	// repositories
	userRepository := postgres.NewUserRepository(postgres.NewUserRepositoryParams{
		Store: postgresDB,
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
	})

	// handlers
	v1 := r.Group("/v1")
	{
		rest.UseUserHandler(rest.NewUserHandlerParams{
			Router:      v1,
			UserService: userService,
		})
		rest.UseAuthHandler(rest.NewAuthHandlerParams{
			Router:      v1,
			AuthService: authService,
		})
	}

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
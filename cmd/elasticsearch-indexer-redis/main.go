package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Stuhub-io/config"
	internalRedis "github.com/Stuhub-io/internal/cache/redis"
	"github.com/Stuhub-io/internal/search/elasticsearch"
	"github.com/Stuhub-io/logger"
)

func main() {
	var env string

	flag.StringVar(&env, "env", "", "Environment Variables filename")
	flag.Parse()

	errC, err := run(env)
	if err != nil {
		log.Fatalf("Couldn't run: %s", err)
	}

	if err := <-errC; err != nil {
		log.Fatalf("Error while running: %s", err)
	}
}

func run(env string) (<-chan error, error) {
	cfg := config.LoadConfig(config.GetConfigLoaders(env))

	logger := logger.NewLogrusLogger()

	redisCache := internalRedis.Must(cfg.RedisUrl, logger)
	redisClient := redisCache.GetClient()

	elasticSearch := elasticsearch.Must(cfg.ElasticSearchURL, logger)

	pageIndexer := elasticsearch.NewPageIndexer(elasticSearch, logger)

	srv := &Server{
		logger:      logger,
		rdb:         redisClient,
		pageIndexer: pageIndexer,
		done:        make(chan struct{}),
	}

	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		<-ctx.Done()

		logger.Info("Shutdown signal received")

		ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer func() {
			redisClient.Close()
			stop()
			cancel()
			close(errC)
		}()

		if err := srv.Shutdown(ctxTimeout); err != nil {
			errC <- err
		}

		logger.Info("Shutdown completed")
	}()

	go func() {
		logger.Info("Listening and serving")

		if err := srv.ListenAndServe(); err != nil {
			errC <- err
		}
	}()

	return errC, nil
}

type Server struct {
	logger      logger.Logger
	rdb         *redis.Client
	pubsub      *redis.PubSub
	pageIndexer *elasticsearch.PageIndexer
	done        chan struct{}
}

func (s *Server) ListenAndServe() error {
	pubsub := s.rdb.PSubscribe(context.Background(), "page.*")

	_, err := pubsub.Receive(context.Background())
	if err != nil {
		return err
	}

	s.pubsub = pubsub

	ch := pubsub.Channel()

	go func() {
		for msg := range ch {
			s.logger.Info("Received message: " + msg.Channel)

			switch msg.Channel {
			case "page.created", "page.updated":
				s.logger.Info("Record saved")
				fmt.Println(msg.Payload)
			}
		}

		s.logger.Info("No more messages to consume. Exiting.")

		s.done <- struct{}{}
	}()

	return nil
}

// Shutdown ...
func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server")

	s.pubsub.Close()

	for {
		select {
		case <-ctx.Done():
			return errors.New("Something wrong")
		case <-s.done:
			return nil
		}
	}
}

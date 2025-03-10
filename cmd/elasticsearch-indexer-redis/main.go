package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	internalRedis "github.com/Stuhub-io/internal/cache/redis"
	"github.com/Stuhub-io/internal/search/elasticsearch"
	"github.com/Stuhub-io/logger"
	"github.com/Stuhub-io/utils/userutils"
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
			s.logger.Info("received message: " + msg.Channel)

			switch msg.Channel {
			case "page.created":
				var page domain.Page

				if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&page); err != nil {
					s.logger.Infof("invalid page payload: %s", err.Error())
					break
				}

				if err := s.pageIndexer.Index(context.Background(), transformDomainToIndexed(page)); err != nil {
					s.logger.Infof("failed to index page: %s", err.Error())
				}

				s.logger.Infof("index page with id %s successfully", page.ID)
			case "page.updated":
				var page domain.Page

				if err := json.NewDecoder(strings.NewReader(msg.Payload)).Decode(&page); err != nil {
					s.logger.Infof("invalid page payload: %s", err.Error())
					break
				}

				if err := s.pageIndexer.Update(context.Background(), transformDomainToIndexed(page)); err != nil {
					s.logger.Infof("failed to update indexed page: %s", err.Error())
				}

				s.logger.Infof("update indexed page with id %s successfully", page.ID)
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

func transformDomainToIndexed(page domain.Page) domain.IndexedPage {
	content := ""
	if page.ViewType == domain.PageViewTypeDoc {
		content = page.Document.JsonContent
	}

	return domain.IndexedPage{
		PkID:           page.PkID,
		ID:             page.ID,
		Name:           page.Name,
		AuthorPkID:     *page.AuthorPkID,
		AuthorFullName: userutils.GetUserFullName(page.Author.FirstName, page.Author.LastName),
		SharedPKIDs:    make([]int64, 0),
		ViewType:       page.ViewType.String(),
		Content:        content,
		CreatedAt:      page.CreatedAt,
		UpdatedAt:      page.UpdatedAt,
		ArchivedAt:     page.ArchivedAt,
	}
}

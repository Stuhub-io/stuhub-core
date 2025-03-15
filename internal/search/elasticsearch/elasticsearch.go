package elasticsearch

import (
	"github.com/Stuhub-io/logger"
	"github.com/elastic/go-elasticsearch/v8"
)

func open(url string, logger logger.Logger) (*elasticsearch.Client, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{url},
	})
	if err != nil {
		logger.Fatalf(err, "failed to open elasticsearch connection")
		return nil, err
	}

	if _, err = es.Ping(); err != nil {
		logger.Fatalf(err, "failed to ping elasticsearch")
		return nil, err
	}

	logger.Info("elasticsearch connected")

	return es, nil
}

func Must(url string, logger logger.Logger) *elasticsearch.Client {
	es, err := open(url, logger)
	if err != nil {
		panic(err)
	}

	return es
}

package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type PageIndexer struct {
	client *elasticsearch.Client
	index  string
}

func NewPageIndexer(client *elasticsearch.Client) *PageIndexer {
	return &PageIndexer{
		client: client,
		index:  "page",
	}
}

func (i *PageIndexer) Index(ctx context.Context) error {
	document := struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}{
		Title:   "nice",
		Content: "asd",
	}
	data, err := json.Marshal(document)
	if err != nil {
		log.Fatalf("Error marshaling the document: %s", err)
		return err
	}

	req := esapi.IndexRequest{
		Index:      i.index,
		DocumentID: "1",
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	resp, err := req.Do(context.Background(), i.client)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
		return err
	}
	defer resp.Body.Close()

	log.Printf("Indexed document %s to index %s\n", resp.String(), i.index)

	return nil
}

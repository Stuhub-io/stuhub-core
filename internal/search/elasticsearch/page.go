package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/logger"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const MaxSearchPagesSize = 5

type PageIndexer struct {
	client *elasticsearch.Client
	logger logger.Logger
	index  string
}

func NewPageIndexer(client *elasticsearch.Client, logger logger.Logger) *PageIndexer {
	return &PageIndexer{
		client: client,
		logger: logger,
		index:  "page",
	}
}

func (i *PageIndexer) Index(ctx context.Context, page domain.IndexedPage) error {
	data, err := json.Marshal(page)
	if err != nil {
		i.logger.Infof("Error marshaling the page: %s", err)
		return err
	}

	req := esapi.IndexRequest{
		Index:      i.index,
		DocumentID: page.ID,
		Body:       bytes.NewReader(data),
		Refresh:    "true",
	}

	resp, err := req.Do(context.Background(), i.client)
	if err != nil {
		i.logger.Infof("Error index the page: %s", err)
		return err
	}
	defer resp.Body.Close()

	i.logger.Infof("Indexed page %s to index %s\n", resp.String(), i.index)

	return nil
}

func (i *PageIndexer) Search(ctx context.Context, args domain.SearchIndexedPageParams) (*[]domain.QuickSearchPage, error) {
	must := make([]interface{}, 0, 4)

	must = append(must, map[string]interface{}{
		"bool": map[string]interface{}{
			"should": []map[string]interface{}{
				{
					"wildcard": map[string]interface{}{
						"name": map[string]interface{}{
							"value":            fmt.Sprintf("*%s*", args.Keyword),
							"case_insensitive": true,
						},
					},
				},
				{
					"wildcard": map[string]interface{}{
						"content": map[string]interface{}{
							"value":            fmt.Sprintf("*%s*", args.Keyword),
							"case_insensitive": true,
						},
					},
				},
			},
			"minimum_should_match": 1,
		},
	})

	if args.ViewType != nil {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"view_type": *args.ViewType,
			},
		})
	}

	if args.AuthorPkID != nil {
		must = append(must, map[string]interface{}{
			"match": map[string]interface{}{
				"author_pkid": *args.AuthorPkID,
			},
		})
	} else {
		must = append(must, map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"author_pkid": args.UserPkID,
						},
					},
					{
						"terms": map[string]interface{}{
							"shared_pkids": []interface{}{args.UserPkID},
						},
					},
				},
				"minimum_should_match": 1,
			},
		})
	}

	query := map[string]interface{}{
		"size": MaxSearchPagesSize,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": must,
			},
		},
		"sort": []map[string]interface{}{
			{"_score": map[string]string{"order": "desc"}},
		},
		// "_source": []string{"id"}, // Specific fields
	}

	var buf bytes.Buffer

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		i.logger.Infof("Error encoding query: %s", err)
		return nil, err
	}

	req := esapi.SearchRequest{
		Index: []string{i.index},
		Body:  &buf,
	}

	resp, err := req.Do(context.Background(), i.client)
	if err != nil {
		i.logger.Infof("Error search indexed pages: %s", err)
		return nil, err
	}
	defer resp.Body.Close()

	var hits struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source domain.IndexedPage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&hits); err != nil {
		return nil, err
	}

	pages := make([]domain.QuickSearchPage, len(hits.Hits.Hits))

	for i, hit := range hits.Hits.Hits {
		if hit.Source.ID == "" {
			pages = make([]domain.QuickSearchPage, 0)
			break
		}
		pages[i].PkID = hit.Source.PkID
		pages[i].ID = hit.Source.ID
		pages[i].Name = hit.Source.Name
		pages[i].AuthorPkID = hit.Source.AuthorPkID
		pages[i].AuthorFullName = hit.Source.AuthorFullName
		pages[i].ViewType = hit.Source.ViewType
		pages[i].UpdatedAt = hit.Source.UpdatedAt
		pages[i].ArchivedAt = hit.Source.ArchivedAt
	}

	return &pages, nil
}

func (i *PageIndexer) Update(ctx context.Context, page domain.IndexedPage) error {
	data, err := json.Marshal(page)
	if err != nil {
		i.logger.Infof("Error marshaling the page: %s", err)
		return err
	}

	req := esapi.UpdateRequest{
		Index:      i.index,
		DocumentID: page.ID,
		Body:       bytes.NewReader([]byte(fmt.Sprintf(`{"doc":%s}`, data))),
	}

	resp, err := req.Do(context.Background(), i.client)
	if err != nil {
		i.logger.Infof("Error update the indexed page: %s", err)
		return err
	}
	defer resp.Body.Close()

	i.logger.Infof("Update indexed page %s to index %s\n", resp.String(), i.index)

	return nil
}

func (i *PageIndexer) Delete(ctx context.Context, pageID string) error {
	req := esapi.DeleteRequest{
		Index:      i.index,
		DocumentID: pageID,
	}

	resp, err := req.Do(context.Background(), i.client)
	if err != nil {
		i.logger.Infof("Error delete the indexed page: %s", err)
		return err
	}
	defer resp.Body.Close()

	i.logger.Infof("Delete indexed page ID %s from index %s\n", pageID, i.index)

	return nil
}

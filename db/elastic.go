package db

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"places/internal/entities"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/some"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/segmentsortorder"
)

type ElasticStore struct {
	ClassicClient *elasticsearch.Client
	TypedClient   *elasticsearch.TypedClient
}

// const (
// 	indexName string = "places"
// 	batch     int    = 250
// )

func ConnectWithElasticSearch() (*ElasticStore, error) {
	es_config := elasticsearch.Config{
		Addresses: []string{
			"http://elasticsearch:9200",
		},
		RetryBackoff: func(attempt int) time.Duration {
			// Exponential backoff
			return time.Duration(attempt) * 100 * time.Millisecond
		},
		MaxRetries: 5,
	}
	newClient, err := elasticsearch.NewClient(es_config)
	if err != nil {
		panic(err)
	}

	newTypedClient, err := elasticsearch.NewTypedClient(es_config)
	if err != nil {
		panic(err)
	}

	return &ElasticStore{newClient, newTypedClient}, err
}

func (e ElasticStore) InsertPlaces(places []entities.Place) (uint64, error) {

	indexName := "places"

	IndexExist, err := e.isIndexExist(indexName)
	if err != nil {
		return 0, err
	}

	if IndexExist {
		err := e.deleteIndex(indexName)
		if err != nil {
			return 0, err
		}
		log.Println("delete old index")
		err = e.createIndex(indexName)
		if err != nil {
			return 0, err
		}
		log.Println("create new index")
	} else {
		err := e.createIndex(indexName)
		if err != nil {
			return 0, err
		}
		log.Println("create new index")
	}

	indexed, err := e.bulkPlaces(places)
	if err != nil {
		return indexed, err
	}
	log.Printf("upload %d places\n", indexed)

	return indexed, nil
}

func (e ElasticStore) isIndexExist(indexN string) (bool, error) {
	res, err := e.ClassicClient.Indices.Exists([]string{indexN})
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return false, errors.New(fmt.Sprintf("Cannot check index exists: %s", err))
	}
	return !res.IsError(), nil
}

func (e ElasticStore) createIndex(indexN string) error {

	tmpFile, err := os.ReadFile("/app/config/schema.json")
	if err != nil {
		return err
	}

	var tmpMap *types.TypeMapping
	err = json.Unmarshal(tmpFile, &tmpMap)
	if err != nil {
		return err
	}

	req := &create.Request{
		Mappings: tmpMap,
		Settings: &types.IndexSettings{
			MaxResultWindow: some.Int(2000),
			Sort: &types.IndexSegmentSort{
				Field: []string{"id"},
				Order: []segmentsortorder.SegmentSortOrder{
					{Name: "asc"},
				},
			},
			RefreshInterval:  "1s",
			NumberOfReplicas: "0",
		},
	}

	res, err := e.TypedClient.Indices.Create(indexN).
		Request(req).
		Do(nil)

	if err != nil {
		return errors.New(fmt.Sprintf("Cannot create index: %s", err))
	}

	if !res.Acknowledged && res.Index != indexN {
		return errors.New(fmt.Sprintf("unexpected error during index creation, got : %#v", res))
	}
	return nil
}

func (e ElasticStore) deleteIndex(indexN string) error {
	res, err := e.ClassicClient.Indices.Delete([]string{indexN}, e.ClassicClient.Indices.Delete.WithIgnoreUnavailable(true))
	if err != nil || res.IsError() {
		log.Fatalf("Cannot delete index: %s", err)
	}
	res.Body.Close()

	return nil
}

func (e ElasticStore) bulkPlaces(places []entities.Place) (uint64, error) {
	log.Printf("Attempting to insert %d places\n", len(places))
	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:      "places",
		Client:     e.ClassicClient,
		NumWorkers: 5,
	})
	if err != nil {
		return 0, err
	}
	defer bulkIndexer.Close(context.Background())

	for _, place := range places {

		jsonPlace, err := json.Marshal(place)
		if err != nil {
			return 0, err
		}

		log.Println(jsonPlace)
		err = bulkIndexer.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: strconv.FormatUint(uint64(place.ID), 10),
				Body:       bytes.NewReader(jsonPlace),
			})
		if err != nil {
			log.Printf("Error adding place to bulk indexer: %v", err)
			continue
		}
	}

	biStats := bulkIndexer.Stats()
	if biStats.NumAdded != uint64(len(places)) {
		return 0, errors.New(fmt.Sprintf("добавлены не все файлы: %d вместо %d", biStats.NumAdded, len(places)))
	}

	stats := bulkIndexer.Stats()
	log.Printf("Bulk indexer stats: %+v", stats)

	return biStats.NumAdded, nil
}

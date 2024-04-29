package storage

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"places/internal/config"
	"places/internal/entities"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type ElasticStore struct {
	ClassicClient *elasticsearch.Client
}

func ConnectWithElasticSearch(logger *slog.Logger) (*ElasticStore, error) {
	es_config := elasticsearch.Config{
		Addresses: []string{
			config.ElasticAddress,
		},
		Username: "elastic",
		Password: "123456",
	}
	classicClient, err := elasticsearch.NewClient(es_config)
	if err != nil {
		logger.Error("Error while creating new client %s", err)
		return nil, err
	}

	res, err := classicClient.Info()
	if err != nil {
		logger.Error("Error while getting client info %s", err)
		return nil, err
	}
	defer res.Body.Close()

	return &ElasticStore{classicClient}, err
}

func (e ElasticStore) InsertPlaces(places []entities.Place, logger *slog.Logger) (uint64, error) {
	IndexExist, err := e.indexExists(config.IndexName)
	if err != nil {
		return 0, err
	}

	if IndexExist {
		err := e.deleteIndex(config.IndexName)
		if err != nil {
			return 0, err
		}
		logger.Info("Deleted old index")
	}

	err = e.createIndex(config.IndexName)
	if err != nil {
		return 0, err
	}
	logger.Info("Created new index")

	indexed, err := e.bulkPlaces(places)
	if err != nil {
		return indexed, err
	}

	_, err = e.ClassicClient.Indices.Refresh(e.ClassicClient.Indices.Refresh.WithIndex(config.IndexName))
	if err != nil {
		logger.Error("Error while refreshing index %s", err)
	}

	return indexed, nil
}

func (e ElasticStore) indexExists(indexN string) (bool, error) {
	res, err := e.ClassicClient.Indices.Exists([]string{indexN})
	if res != nil {
		defer res.Body.Close()
	}

	if err != nil {
		return false, err
	}
	return !res.IsError(), nil
}

func (e ElasticStore) createIndex(indexN string) error {

	tmpFile, err := os.ReadFile(config.Schema)
	if err != nil {
		return err
	}

	var mapping *types.TypeMapping
	err = json.Unmarshal(tmpFile, &mapping)
	if err != nil {
		return err
	}

	_, err = e.ClassicClient.Indices.Create(indexN)

	if err != nil {
		return err
	}

	return nil
}

func (e ElasticStore) deleteIndex(indexN string) error {
	res, err := e.ClassicClient.Indices.Delete([]string{indexN}, e.ClassicClient.Indices.Delete.WithIgnoreUnavailable(true))
	if err != nil || res.IsError() {
		return err
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
		return 0, err
	}

	stats := bulkIndexer.Stats()
	log.Printf("Bulk indexer stats: %+v", stats)

	return biStats.NumAdded, nil
}

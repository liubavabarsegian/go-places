package db

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

func ConnectWithElasticSearch() (*ElasticStore, error) {
	es_config := elasticsearch.Config{
		Addresses: []string{
			config.ElasticAddress,
		},
		Username: "elastic",
		Password: "123456",
	}
	classicClient, err := elasticsearch.NewClient(es_config)
	if err != nil {
		panic(err)
	}

	res, err := classicClient.Info()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return &ElasticStore{classicClient}, err
}

func (e ElasticStore) InsertPlaces(places []entities.Place) (uint64, error) {
	IndexExist, err := e.indexExists(config.IndexName)
	if err != nil {
		return 0, err
	}

	if IndexExist {
		err := e.deleteIndex(config.IndexName)
		if err != nil {
			return 0, err
		}
		log.Println("deletde old index")
	}

	err = e.createIndex(config.IndexName)
	if err != nil {
		return 0, err
	}
	log.Println("created new index")

	indexed, err := e.bulkPlaces(places)
	if err != nil {
		return indexed, err
	}

	_, err = e.ClassicClient.Indices.Refresh(e.ClassicClient.Indices.Refresh.WithIndex(config.IndexName))
	if err != nil {
		log.Println(err)
	}

	return indexed, nil
}

func (e ElasticStore) indexExists(indexN string) (bool, error) {
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
		return errors.New(fmt.Sprintf("Cannot create index: %s", err))
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

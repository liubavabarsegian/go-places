package db

import (
	"context"
	"errors"
	"fmt"
	"places/internal/entities"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type Store interface {
	GetPlaces(limit int, offset int) ([]entities.Place, int, error)
}

func (s ElasticStore) GetPlaces(limit int, offset int) ([]entities.Place, int, error) {
	fmt.Println("AAA")
	if limit <= 0 {
		return nil, 0, errors.New("limit must be > 0")
	}
	if offset < 0 {
		return nil, 0, errors.New("offset must be >= 0")
	}

	fmt.Println("BBBB")
	// query := map[string]interface{}{
	// 	"query": map[string]interface{}{
	// 		"match_all": map[string]interface{}{},
	// 	},
	// }

	// var buf bytes.Buffer
	// if err := json.NewEncoder(&buf).Encode(query); err != nil {
	// 	return nil, 0, err
	// }

	qry := &search.Request{
		Size: &limit,
		From: &offset,
		Query: &types.Query{
			MatchAll: &types.MatchAllQuery{Boost: nil, QueryName_: nil},
		},
		TrackTotalHits: true,
	}
	// log.Printf("Sending query to Elasticsearch: %s", buf.String())

	res, err := s.TypedClient.Search().Index("places").Request(qry).Do(context.Background())
	// s.ClassicClient.Search.WithIndex("places"),
	// // s.ClassicClient.Search.WithFrom(offset),
	// s.ClassicClient.Search.WithSize(3),
	// s.ClassicClient.Search.WithPretty(),
	// s.ClassicClient.Search.WithBody(&buf),
	// s.ClassicClient.Search.WithTrackTotalHits(true),
	// s.ClassicClient.Search.WithContext(context.Background()),

	fmt.Println(res)
	if err != nil {
		return nil, 0, err
	}
	// if res.IsError() {
	// 	// Log the error message for debugging
	// 	log.Printf("Error in Elasticsearch operation: %s", res.Status())
	// 	return nil, 0, errors.New("Elasticsearch operation error")
	// }

	// defer res.Body.Close()

	// var response map[string]interface{}
	// if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
	// 	return nil, 0, err
	// }

	// // Log the raw response from Elasticsearch
	// log.Printf("Raw response from Elasticsearch: %+v", response)

	// hits, found := response["hits"].(map[string]interface{})
	// if !found {
	// 	return nil, 0, errors.New("No hits found in response")
	// }

	// totalHits, _ := hits["total"].(map[string]interface{})["value"].(float64)
	// var places []entities.Place

	// if int(totalHits) > 0 {
	// 	for _, hit := range hits["hits"].([]interface{}) {
	// 		source := hit.(map[string]interface{})["_source"]
	// 		place := entities.Place{} // You need to define your Place type
	// 		placeBytes, err := json.Marshal(source)
	// 		if err != nil {
	// 			return nil, 0, err
	// 		}
	// 		if err := json.Unmarshal(placeBytes, &place); err != nil {
	// 			return nil, 0, err
	// 		}
	// 		places = append(places, place)
	// 	}
	// }

	// return places, int(totalHits), nil
	return nil, 0, nil
}

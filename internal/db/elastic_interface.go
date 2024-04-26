package db

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"places/internal/config"
	"places/internal/entities"
)

type Store interface {
	GetPlaces(limit int, offset int) ([]entities.Place, int, error)
}

type searchRequestParams struct {
	Took    float64 `json:"took"`
	Timeout bool    `json:"timed_out"`
	Shards  struct {
		Total      int64 `json:"total"`
		Successful int64 `json:"successful"`
		Skipped    int64 `json:"skipped"`
		Failed     int64 `json:"failed"`
	} `json:"_shards"`
	Hits struct {
		Total struct {
			Value    int64  `json:"value"`
			Relation string `json:"relation"`
		} `json:"total"`
		MaxScore float64 `json:"max_score"`
		Hits     []struct {
			Index  string          `json:"_index"`
			Id     string          `json:"_id"`
			Score  float64         `json:"_score"`
			Source *entities.Place `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (e ElasticStore) GetPlaces(limit int, offset int) ([]entities.Place, int, error) {
	if limit <= 0 {
		return nil, 0, errors.New("limit must be > 0")
	}
	if offset < 0 {
		return nil, 0, errors.New("offset must be >= 0")
	}

	response, err := e.ClassicClient.Search(
		e.ClassicClient.Search.WithContext(context.Background()),
		e.ClassicClient.Search.WithIndex(config.IndexName),
		e.ClassicClient.Search.WithFrom(offset),
		e.ClassicClient.Search.WithSize(10),
		e.ClassicClient.Search.WithSort("ID:asc"),
		e.ClassicClient.Search.WithTrackTotalHits(true),
		e.ClassicClient.Search.WithPretty(),
	)

	var respParams searchRequestParams
	err = json.NewDecoder(response.Body).Decode(&respParams)
	if err != nil {
		return nil, 0, err
	}

	var places []entities.Place
	if respParams.Hits.Total.Value > 0 {
		for _, hit := range respParams.Hits.Hits {
			if hit.Source == nil {
				log.Printf("hit with %s have nil Source", hit.Id)
				continue
			}
			places = append(places, *hit.Source)
		}
	}

	return places, int(respParams.Hits.Total.Value), nil
}

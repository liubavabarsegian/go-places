package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"places/internal/config"
	"places/internal/entities"
	"strings"
)

type Store interface {
	GetPlaces(limit int, offset int) ([]entities.Place, int, error)
}

type searchResponse struct {
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

func (e ElasticStore) GetPlaces(limit int, offset int, logger *slog.Logger) ([]entities.Place, int, error) {
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
		e.ClassicClient.Search.WithTrackTotalHits(true),
		e.ClassicClient.Search.WithPretty(),
	)

	if err != nil {
		logger.Error("Error while searching places")
		return nil, 0, err
	}

	var responseParams searchResponse
	err = json.NewDecoder(response.Body).Decode(&responseParams)
	if err != nil {
		return nil, 0, err
	}

	var places []entities.Place
	if responseParams.Hits.Total.Value > 0 {
		for _, hit := range responseParams.Hits.Hits {
			if hit.Source == nil {
				logger.Info(fmt.Sprintf("hit with %s have nil Source", hit.Id))
				continue
			}
			places = append(places, *hit.Source)
		}
	}

	return places, int(responseParams.Hits.Total.Value), nil
}

func (e ElasticStore) GetClosestPlaces(longitude float64, latitude float64, logger *slog.Logger) ([]entities.Place, int, error) {
	query := map[string]interface{}{
		"sort": []map[string]interface{}{
			{
				"_geo_distance": map[string]interface{}{
					"location": map[string]interface{}{
						"lat": latitude,
						"lon": longitude,
					},
					"order":           "asc",
					"unit":            "km",
					"mode":            "min",
					"distance_type":   "arc",
					"ignore_unmapped": true,
				},
			},
		},
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		logger.Error("Error ocurred while marshalling query")
		return nil, 0, err
	}

	response, err := e.ClassicClient.Search(
		e.ClassicClient.Search.WithContext(context.Background()),
		e.ClassicClient.Search.WithIndex(config.IndexName),
		e.ClassicClient.Search.WithSize(3),
		e.ClassicClient.Search.WithBody(strings.NewReader(string(queryBytes))),
		e.ClassicClient.Search.WithPretty(),
	)

	if err != nil {
		logger.Error("Error while searching the closest places")
		return nil, 0, err
	}

	var responseParams searchResponse
	err = json.NewDecoder(response.Body).Decode(&responseParams)
	if err != nil {
		return nil, 0, err
	}

	var places []entities.Place
	if responseParams.Hits.Total.Value > 0 {
		for _, hit := range responseParams.Hits.Hits {
			if hit.Source == nil {
				logger.Info(fmt.Sprintf("hit with %s have nil Source", hit.Id))
				continue
			}
			places = append(places, *hit.Source)
		}
	}

	return places, int(responseParams.Hits.Total.Value), nil
}

package db

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"places/config"
	"places/internal/entities"
	"places/internal/serializers"
)

type Store interface {
	GetPlaces(limit int, offset int) ([]entities.Place, int, error)
}

func (s ElasticStore) GetPlaces(limit int, offset int) ([]entities.Place, int, error) {
	if limit <= 0 {
		return nil, 0, errors.New("limit mast be > 0")
	}
	if offset < 0 {
		return nil, 0, errors.New("offset mast be >= 0")
	}

	var result []entities.Place

	for i := offset; i < 100; i += limit {
		res, err := s.ClassicClient.Search(
			s.ClassicClient.Search.WithContext(context.Background()),
			s.ClassicClient.Search.WithIndex(config.IndexName),
			s.ClassicClient.Search.WithFrom(i),
			s.ClassicClient.Search.WithSize(limit),
			s.ClassicClient.Search.WithTrackTotalHits(true),
			s.ClassicClient.Search.WithPretty(),
		)
		if err != nil || res.IsError() {
			return nil, 0, err
		}

		defer res.Body.Close()

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, 0, err
		}

		var responses serializers.Response
		if err := json.Unmarshal(bodyBytes, &responses); err != nil {
			return nil, 0, err
		}

		totalHits := responses.Hits.Total.Value
		var places []entities.Place

		for _, response := range responses.Hits.Hits {
			places = append(places, response.Source)
		}

		return places, totalHits, nil

		// users := res.PLaces
		// result = append(result, res.Body)
	}

	return result, len(result), nil

	// req := &search.Request{
	// 	From: &offset,
	// 	Size: &limit,
	// }

	// res, err := s.TypedClient.Search().
	// 	Index(config.IndexName).
	// 	Request(req).
	// 	Do(context.Background())

	// if err != nil {
	// 	return nil, 0, err
	// }

	// rawPlaces := res.Hits.Hits

	// places := make([]entities.Place, len(rawPlaces))
	// for i := 0; i < len(res.Hits.Hits); i++ {
	// 	err := json.Unmarshal(rawPlaces[i].Source_, &places[i])
	// 	if err != nil {
	// 		return nil, 0, err
	// 	}
	// }

	// return places, int(res.Hits.Total.Value), nil

	// res, err := p.elastic.client.Search()
	// req := esapi.SearchRequest{
	// 	Index: []string{config.IndexName},
	// 	From:  &offset,
	// 	Size:  &limit,
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 100)
	// defer cancel()

	// res, err := req.Do(ctx, s.ClassicClient)
	// if err != nil {
	// 	return nil, 0, fmt.Errorf("list all: request: %w", err)
	// }
	// defer res.Body.Close()

	// if res.IsError() {
	// 	return nil, 0, fmt.Errorf("list all: response: %s", res.String())
	// }

	// var body postsResponse
	// if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
	// 	return nil, 0, fmt.Errorf("list all: decode: %w", err)
	// }

	// posts := make([]*storage.Post, len(body.Hits.Hits))
	// for i, v := range body.Hits.Hits {
	// 	posts[i] = v.Source
	// }
	//
	// return body.Hits.Total.Value, posts, nil
}

// type postsResponse struct {
// 	Hits struct {
// 		Total struct {
// 			Value int `json:"value"`
// 		} `json:"total"`
// 		Hits []struct {
// 			Source *storage.Post `json:"_source"`
// 		} `json:"hits"`
// 	} `json:"hits"`
// }

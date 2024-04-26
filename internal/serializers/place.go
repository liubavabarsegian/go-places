package serializers

import "places/internal/entities"

type GetPlacesResponse struct {
	Name     string           `json:"name"`
	Total    int              `json:"total"`
	Places   []entities.Place `json:"places"`
	PrevPage int              `json:"prev_page"`
	NextPage int              `json:"next_page"`
	LastPage int              `json:"last_page"`
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

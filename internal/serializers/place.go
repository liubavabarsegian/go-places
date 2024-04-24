package serializers

import "places/internal/entities"

type GetPlacesResponse struct {
}

type Response struct {
	Hits Hits `json:"hits"`
}

type Hits struct {
	Total Total `json:"total"`
	Hits  []Hit `json:"hits"`
}

type Total struct {
	Value int `json:"value"`
}

type Hit struct {
	Source entities.Place `json:"_source"`
}

type PlacePageData struct {
	Places  []entities.Place
	Total   int
	Name    string
	Address string
	Phone   string
}

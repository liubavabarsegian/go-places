package models

import (
	// "PlacesApp/internal"
	"PlacesApp/internal"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"strconv"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	esapi "github.com/elastic/go-elasticsearch/v8/esapi"
	"go.opentelemetry.io/otel"
)

const (
	ErrInvalidPlaceID = "invalid place ID"
)

type Place struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	Address   string   `json:"address"`
	Phone     string   `json:"phone"`
	Location  GeoPoint `json:"location"`
	Client    *elasticsearch.Client
	IndexName string
}

type GeoPoint struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

var Places []Place

func (place *Place) Validate() error {
	// place_id, _ := strconv.Atoi(place.ID)
	if place.ID < 1 {
		return errors.New(ErrInvalidPlaceID)
	}
	return nil
}

func NewPlace(client *elasticsearch.Client) *Place {
	return &Place{
		Client:    client,
		IndexName: "places",
	}
}

// Index creates or updates a place in an index.
func (t *Place) Index(ctx context.Context, place Place) error {
	tracer := otel.Tracer("PlacesApp")
	ctx, span := tracer.Start(ctx, "Place.Index")
	defer span.End()

	body := Place{
		ID:       place.ID,
		Name:     place.Name,
		Address:  place.Address,
		Phone:    place.Phone,
		Location: GeoPoint(place.Location),
	}

	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(body) // XXX: error omitted

	req := esapi.IndexRequest{
		Index:      t.IndexName,
		Body:       &buf,
		DocumentID: strconv.Itoa(place.ID),
		Refresh:    "true",
	}

	resp, err := req.Do(ctx, t.Client)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "IndexRequest.Do")
	}
	defer resp.Body.Close()

	io.Copy(ioutil.Discard, resp.Body)

	return nil
}

package repository

import (
	"PlacesApp/internal/db/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"PlacesApp/internal"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	esapi "github.com/elastic/go-elasticsearch/v8/esapi"
	"go.opentelemetry.io/otel"
)

type Place struct {
	client *elasticsearch.Client
	index  string
}

type indexedPlace struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location GeoPoint `json:"location"`
}

type GeoPoint struct {
	Longitude float64 `json:"lon"`
	Latitide  float64 `json:"lat"`
}

func NewPlace(client *elasticsearch.Client) *Place {
	return &Place{
		client: client,
		index:  "places",
	}
}

// Index creates or updates a task in an index.
func (t *Place) Index(ctx context.Context, place models.Place) error {
	// Obtain a Tracer from the global TracerProvider
	tracer := otel.Tracer("PlacesApp")

	// Start a new span
	ctx, span := tracer.Start(ctx, "Place.Index")
	defer span.End()

	body := indexedPlace{
		ID:       place.ID,
		Name:     place.Name,
		Address:  place.Address,
		Phone:    place.Phone,
		Location: GeoPoint(place.Location),
	}

	var buf bytes.Buffer

	_ = json.NewEncoder(&buf).Encode(body) // XXX: error omitted

	req := esapi.IndexRequest{
		Index:      t.index,
		Body:       &buf,
		DocumentID: place.ID,
		Refresh:    "true",
	}

	resp, err := req.Do(ctx, t.client)
	if err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "IndexRequest.Do")
	}
	defer resp.Body.Close()

	// if resp.IsError() {
	// 	return internal.NewErrorf(internal.ErrorCodeUnknown, "IndexRequest.Do %s", resp.StatusCode)
	// }

	io.Copy(ioutil.Discard, resp.Body)

	return nil
}

func ParseBody(r *http.Request, x interface{}) {
	if body, err := ioutil.ReadAll(r.Body); err == nil {
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}

func Index(w http.ResponseWriter, r *http.Request, ctx context.Context) error {
	place := &models.Place{}
	ParseBody(r, place)
	// Obtain a Tracer from the global TracerProvider
	tracer := otel.Tracer("PlacesApp")

	// Start a new span
	ctx, span := tracer.Start(ctx, "Place.Index")
	defer span.End()

	body := indexedPlace{
		ID:       place.ID,
		Name:     place.Name,
		Address:  place.Address,
		Phone:    place.Phone,
		Location: GeoPoint(place.Location),
	}

	var buf bytes.Buffer

	_ = json.NewEncoder(&buf).Encode(body) // XXX: error omitted

	req := esapi.IndexRequest{
		Index:      "t.index",
		Body:       &buf,
		DocumentID: place.ID,
		Refresh:    "true",
	}

	fmt.Println(req)
	return nil
}

// type TaskRepository interface {
// 	Create(ctx context.Context, description string, priority internal.Priority, dates internal.Dates) (internal.Task, error)
// 	Find(ctx context.Context, id string) (internal.Task, error)
// 	Update(ctx context.Context, id string, description string, priority internal.Priority, dates internal.Dates, isDone bool) error
// }

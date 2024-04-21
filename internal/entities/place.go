package entities

import (
	"errors"
)

const (
	ErrInvalidPlaceID = "invalid place ID"
)

type Place struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Address  string   `json:"address"`
	Phone    string   `json:"phone"`
	Location GeoPoint `json:"location"`
	// Client    *elasticsearch.Client
	// IndexName string
}

type GeoPoint struct {
	Longitude float64 `json:"lon"`
	Latitude  float64 `json:"lat"`
}

var Places []Place

func (place *Place) Validate() error {
	if place.ID < 1 {
		return errors.New(ErrInvalidPlaceID)
	}
	return nil
}

// func NewPlace(client *elasticsearch.Client) *Place {
// 	return &Place{
// 		Client:    client,
// 		IndexName: "places",
// 	}
// }

// type TaskRepository interface {
// 	Create(ctx context.Context, description string, priority internal.Priority, dates internal.Dates) (internal.Task, error)
// 	Find(ctx context.Context, id string) (internal.Task, error)
// 	Update(ctx context.Context, id string, description string, priority internal.Priority, dates internal.Dates, isDone bool) error
// }

// Index creates or updates a place in an index.
// func (t *Place) Index(ctx context.Context, place Place) error {
// 	tracer := otel.Tracer("places")
// 	ctx, span := tracer.Start(ctx, "Place.Index")
// 	defer span.End()

// 	body := Place{
// 		ID:       place.ID,
// 		Name:     place.Name,
// 		Address:  place.Address,
// 		Phone:    place.Phone,
// 		Location: GeoPoint(place.Location),
// 	}

// 	var buf bytes.Buffer
// 	json.NewEncoder(&buf).Encode(body)

// 	// req := esapi.IndexRequest{
// 	// 	Index:      t.IndexName,
// 	// 	Body:       &buf,
// 	// 	DocumentID: strconv.Itoa(place.ID),
// 	// 	Refresh:    "true",
// 	// }

// 	var mapping *types.TypeMapping
// 	file, err := os.Open("config/schema.json")
// 	defer file.Close()
// 	byteValue, _ := io.ReadAll(file)
// 	err = json.Unmarshal(byteValue, &mapping)
// 	if err != nil {
// 		return err
// 	}
// 	req := &create.Request{
// 		Mappings: mapping,
// 		Settings: &types.IndexSettings{
// 			MaxResultWindow: some.Int(20000),
// 			Sort: &types.IndexSegmentSort{
// 				Field: []string{"id"},
// 				// Order: []
// 			},
// 		},
// 	}

// 	// resp, err := req.Do(ctx, t.Client)
// 	// if err != nil {
// 	// 	return internal.WrapErrorf(err, internal.ErrorCodeUnknown, "IndexRequest.Do")
// 	// }
// 	// defer resp.Body.Close()

// 	// io.Copy(ioutil.Discard, resp.Body)

// 	// return nil
// 	res, err := elasticsearch.TypedClient.Indices.Create(indexN).
// 		Request(req).
// 		Do(nil)

// 	if err != nil {
// 		return errors.New(fmt.Sprintf("Cannot create index: %s", err))
// 	}

// 	if !res.Acknowledged && res.Index != indexN {
// 		return errors.New(fmt.Sprintf("unexpected error during index creation, got : %#v", res))
// 	}
// 	return nil
// }

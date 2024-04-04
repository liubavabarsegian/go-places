package models

import (
	"errors"
	"strconv"
)

const (
	ErrInvalidPlaceID = "invalid place ID"
)

type Place struct {
	ID       string
	Name     string
	Address  string
	Phone    string
	Location GeoPoint
}

type GeoPoint struct {
	Longitude float64
	Latitide  float64
}

var Places []Place

func (place *Place) Validate() error {
	place_id, _ := strconv.Atoi(place.ID)
	if place_id < 1 {
		return errors.New(ErrInvalidPlaceID)
	}
	return nil
}

// func main() {

// 	client, err := elasticsearch.NewDefaultClient()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	index := "products"
// 	mapping := `
//     {
//       "settings": {
//         "number_of_shards": 1
//       },
//       "mappings": {
//         "properties": {
//           "field1": {
//             "type": "text"
//           }
//         }
//       }
//     }`

// 	res, err := client.Indices.Create(
// 		index,
// 		client.Indices.Create.WithBody(strings.NewReader(mapping)),
// 	)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Println(res)
// }

package repository

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"places/internal/entities"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/pkg/errors"
)

/*************************************************/
/*********** Internal for the Context ************/
/*************************************************/

const (
	indexName string = "places"
	batch     int    = 250
)

type contextKey struct {
	Key int
}

var PlacesKey contextKey = contextKey{Key: 1}
var ClientKey contextKey = contextKey{Key: 2}

func ParsePlacesFromCsv(path string) ([]entities.Place, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var places []entities.Place

	fmt.Println("what")

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '\t'
		return r
	})

	fmt.Println("is")
	err = gocsv.Unmarshal(file, &places)
	if err != nil {
		return nil, err
	}

	fmt.Println("happening")
	return places, err
}

func InsertPlacesIntoElastic(es *elasticsearch.Client, places []entities.Place) error {
	// Определение схемы для маппинга
	// mapping := `{
	// 	"properties": {
	// 	  "id": {
	// 		"type": "long"
	// 	  },
	// 	  "name": {
	// 		"type": "text"
	// 	  },
	// 	  "address": {
	// 		"type": "text"
	// 	  },
	// 	  "phone": {
	// 		"type": "text"
	// 	  },
	// 	  "location": {
	// 		"type": "geo_point"
	// 	  }
	// 	}
	//   }`

	// Добавление маппинга в индекс
	// Добавление маппинга перед началом индексации данных
	// if err := addMapping(es, "places", mapping); err != nil {
	// 	log.Fatalf("Error adding mapping: %s", err)
	// }

	// res, err := es.Indices.PutMapping(
	// 	[]string{"places"},
	// 	strings.NewReader(mapping),
	// 	es.Indices.PutMapping.WithContext(context.Background()),
	// 	// es.Indices.PutMapping.WithIncludeTypeName(true),
	// )
	// if err != nil {
	// 	log.Fatalf("Error adding mapping: %s", err)
	// }
	// defer res.Body.Close()

	fmt.Print("→ Sending batch ")
	var buf bytes.Buffer
	for i, place := range places {
		meta := []byte(fmt.Sprintf(`{ "index" : { "_id" : "%d" } }%s`, place.ID, "\n"))

		data, err := json.Marshal(place)
		if err != nil {
			return errors.Wrap(err, "Cannot encode place "+strconv.Itoa(place.ID))
		}
		data = append(data, "\n"...)

		buf.Grow(len(meta) + len(data))
		buf.Write(meta)
		buf.Write(data)

		if (i+1)%batch == 0 || i == len(places)-1 {
			fmt.Printf("[%d/%d] ", i+1, len(places))

			// var err error
			for i := 0; i < 10; i++ { // Retry up to 10 times
				res, err := es.Bulk(bytes.NewReader(buf.Bytes()), es.Bulk.WithIndex(indexName))
				if err == nil && res != nil {
					// Success, break the loop
					break
				}
				log.Printf("Attempt %d: failure indexing batch %d, retrying...\n", i+1, place.ID)
				time.Sleep(5 * time.Second) // Wait for 5 seconds before retrying
				if res != nil && res.IsError() {
					var raw map[string]interface{}
					if err := json.NewDecoder(res.Body).Decode(&raw); err != nil {
						return errors.Wrap(err, "failure to to parse response body")
					} else {
						log.Printf("  Error: [%d] %s: %s",
							res.StatusCode,
							raw["error"].(map[string]interface{})["type"],
							raw["error"].(map[string]interface{})["reason"],
						)
					}
					res.Body.Close()
				}

				buf.Reset()
			}
		}
	}

	return nil
}

// func addMapping(es *elasticsearch.Client, indexName, mapping string) error {
// 	var buf bytes.Buffer
// 	buf.WriteString(mapping)

// 	req := esapi.IndicesPutMappingRequest{
// 		Index: []string{indexName},
// 		Body:  &buf,
// 	}

// 	res, err := req.Do(context.Background(), es)
// 	if err != nil {
// 		return err
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		var e map[string]interface{}
// 		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
// 			return err
// 		}
// 		return errors.New(fmt.Sprintf("Error adding mapping: %s", e["error"].(map[string]interface{})["reason"]))
// 	}

// 	return nil
// }

package repository

import (
	"encoding/csv"
	"io"
	"os"
	"places/internal/entities"

	"github.com/gocarina/gocsv"
)

func ParsePlacesFromCsv(path string) ([]entities.Place, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var places []entities.Place

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '\t'
		return r
	})

	err = gocsv.Unmarshal(file, &places)
	if err != nil {
		return nil, err
	}

	return places, err
}

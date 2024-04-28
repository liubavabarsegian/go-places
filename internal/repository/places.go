package repository

import (
	"encoding/csv"
	"log/slog"
	"os"
	"places/internal/entities"
	"strconv"
)

func ParsePlacesFromCsv(path string, logger *slog.Logger) ([]entities.Place, error) {
	var data []entities.Place

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t'
	info, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range info[1:] {
		id, err := strconv.Atoi(record[0])
		if err != nil {
			logger.Error("id %s converting error: %s", record[0], err)
			continue
		}

		lon, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			logger.Error("id %d latitude converting error: %s", id, err)
			continue
		}

		lat, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			logger.Error("id %d longitude converting error: %s", id, err)
			continue
		}

		data = append(data, entities.Place{
			ID:      id + 1,
			Name:    record[1],
			Address: record[2],
			Phone:   record[3],
			Location: entities.GeoPoint{
				Longitude: lon,
				Latitude:  lat}})
	}

	return data, nil
}

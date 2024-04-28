package entities

type Place struct {
	ID       int      `csv:"ID" json:"id"`
	Name     string   `csv:"Name" json:"name"`
	Address  string   `csv:"Address" json:"address"`
	Phone    string   `csv:"Phone" json:"phone"`
	Location GeoPoint `json:"location"`
}

type GeoPoint struct {
	Longitude float64 `csv:"Longitude" json:"lon"`
	Latitude  float64 `csv:"Latitude" json:"lat"`
}

var Places []Place

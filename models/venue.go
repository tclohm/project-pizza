package models

type Venue struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Address string `json:"address"`
}
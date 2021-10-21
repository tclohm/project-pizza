package models

type Venue struct {
	Base
	Name string `json:"name"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon`
	Address string `json:"address"`
	Pizzas []Pizza `json:"pizzas"`
}
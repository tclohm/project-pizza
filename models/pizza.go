package models

type Pizza struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Venue Venue `json:"venue"`
	Taste Taste
}
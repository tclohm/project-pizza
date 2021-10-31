package models

type VenuePizza struct {
	ID uint `json:"id"`
	VenueId uint `json:"venue_id"`
	PizzaId uint `json:"pizza_id"`
}
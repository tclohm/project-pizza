package data

type VenuePizza struct {
	ID int64 `json:"id"`
	VenueId int64 `json:"venue_id"`
	PizzaId int64 `json:"pizza_id"`
}
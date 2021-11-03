package models

type Taste struct {
	ID uint `json:"id"`
	Cheesiness float32 `json:"cheesiness"`
	Flavor float32 `json:"flavor"`
	Sauciness float32 `json:"sauciness"`
	Saltiness float32 `json:"saltiness"`
	Charness float32 `json:"charness"`
}
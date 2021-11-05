package models

type Pizza struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Style string `json:"style"`
	Description string `json:"description"`
	Cheesiness float32 `json:"cheesiness"`
	Flavor float32 `json:"flavor"`
	Sauciness float32 `json:"sauciness"`
	Saltiness float32 `json:"saltiness"`
	Charness float32 `json:"charness"`
	ImageId uint `json:"image_id"`
}
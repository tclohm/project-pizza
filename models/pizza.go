package models

type Pizza struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Style string `json:"style"`
	Description string `json:"description"`
	TasteId uint `json:"taste_id"`
	ImageId uint `json:"image_id"`
}
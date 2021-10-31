package models

type Image struct {
	ID uint `json:"id"`
	Filename string `json:"filename"`
	Content_type string `json:"content_type"`
	Location string `json:"location"`
}
package models

type Image struct {
	Filename string `json:"filename"`
	Content_type string `json:"content_type"`
	Location string `json:"location"`
}
package data

type Image struct {
	ID int64 `json:"id"`
	Filename string `json:"filename"`
	Content_type string `json:"content_type"`
	Location string `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}
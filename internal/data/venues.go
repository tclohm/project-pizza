package data

type Venue struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Address string `json:"address"`
}
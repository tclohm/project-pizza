package models

type Taste struct {
	ID uint `json:"id"`
	Cheesiness int `json:"cheesiness"`
	Flavor int `json:"flavor"`
	Sauciness int `json:"sauciness"`
	Saltiness int `json:"saltiness"`
	Charnes int `json:"charness"`
}
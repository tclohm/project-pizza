package models

import (
	"fmt"
)

type Taste struct {
	ID string `json:"Id"`
	Cheesiness int `json:"Cheesiness"`
	Flavor int `json:"Flavor"`
	Sauciness int `json:"Sauciness"`
	Saltiness int `json:"Saltiness"`
	Char int `json:"Char"`
}
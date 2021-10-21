package models

import (
	"fmt"
	"github.com/google/uuid"
)

type Venue struct {
	id ID
	Name string
	Lat float64
	Lon float64
	Address string
	Pizzas []Pizza
}
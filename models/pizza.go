package models

import (
	"fmt"
)

type Pizza struct {
	Name string
	Description string
	Venue Venue
	Taste Taste
}
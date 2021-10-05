package models

import (
	"fmt"
)

type Pizza struct {
	Name string
	Shape string
	Description string
	images []string
	venues []Venue
}
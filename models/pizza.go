package models

import (
	"fmt"
	"github.com/google/uuid"
)

type Pizza struct {
	ID id
	Name string
	Description string
	Venue Venue
	Taste Taste
}
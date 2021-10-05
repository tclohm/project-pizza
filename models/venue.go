package models

import (
	"fmt"
)

type Venue struct {
	Name string
	Capacity string
	Tags []string
	Service []int
	Pizza []Pizza
	Checkout int
}
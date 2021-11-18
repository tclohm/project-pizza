package data

import (
	"time"

	"github.com/tclohm/project-pizza/internal/validator"
)

type Pizza struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Style string `json:"style"`
	Description string `json:"description"`
	Cheesiness float32 `json:"cheesiness"`
	Flavor float32 `json:"flavor"`
	Sauciness float32 `json:"sauciness"`
	Saltiness float32 `json:"saltiness"`
	Charness float32 `json:"charness"`
	ImageFilename string `json:"filename"`
	ImageContentType string `json:"content_type"`
	ImageLocation string `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}

func ValidatePizza(v *validator.Validator, pizza *Pizza) {
	v.Check(pizza.Name != "", "name", "must be provided")
	v.Check(len(pizza.Name) < 500, "name", "must not be more than 500 bytes long")

	v.Check(pizza.Style != "", "style", "must be provided")
	v.Check(len(pizza.Style) < 500, "style", "must not be more than 500 bytes long")

	v.Check(pizza.Description != "", "description", "must be provided")
	v.Check(len(pizza.Description) < 500, "description", "must not be more than 500 bytes long")

	v.Check(pizza.Cheesiness >= 0, "cheesiness", "must be greater than or equal to 0")
	v.Check(pizza.Cheesiness <= 5, "cheesiness", "must be less than or equal to 5")

	v.Check(pizza.Flavor >= 0, "flavor", "must be greater than or equal to 0")
	v.Check(pizza.Flavor <= 5, "flavor", "must be less than or equal to 5")

	v.Check(pizza.Sauciness >= 0, "sauciness", "must be greater than or equal to 0")
	v.Check(pizza.Sauciness <= 5, "sauciness", "must be less than or equal to 5")

	v.Check(pizza.Saltiness >= 0, "saltiness", "must be greater than or equal to 0")
	v.Check(pizza.Saltiness <= 5, "saltiness", "must be less than or equal to 5")

	v.Check(pizza.Charness >= 0, "charness", "must be greater than or equal to 0")
	v.Check(pizza.Charness <= 5, "charness", "must be less than or equal to 5")
}
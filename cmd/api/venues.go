package main

import (
	"fmt"
	"net/http"

	"github.com/tclohm/project-pizza/internal/data"
	"github.com/tclohm/project-pizza/internal/validator"

	_"github.com/gorilla/mux"
)

func (app *application) createVenueHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name 	string `json:"name"`
		Lat 	float64 `json:"lat"`
		Lon 	float64 `json:"lon"`
		Address string `json:"address"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	venue := &data.Venue{
		Name: input.Name,
		Lat: input.Lat,
		Lon: input.Lon,
		Address: input.Address,
	}

	v := validator.New()

	if data.ValidateVenue(v, venue); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Venues.Insert(venue)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/venues/%d", venue.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"venue": venue}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
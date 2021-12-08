package main

import (
	"fmt"
	"net/http"
	"strconv"
	"errors"

	"github.com/tclohm/project-pizza/internal/data"
	"github.com/tclohm/project-pizza/internal/validator"

	"github.com/gorilla/mux"
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

func (app *application) updateVenueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	venue, err := app.models.Venues.Get(n)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name 	*string `json:"name"`
		Lat 	*float64 `json:"lat"`
		Lon 	*float64 `json:"lon"`
		Address *string `json:"address"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		venue.Name = *input.Name
	}

	if input.Lat != nil {
		venue.Lat = *input.Lat
	}

	if input.Lon != nil {
		venue.Lon = *input.Lon
	}

	if input.Address != nil {
		venue.Address = *input.Address
	}

	v := validator.New()

	if data.ValidateVenue(v, venue); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Venues.Update(venue)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"venue": venue}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteVenueHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Venues.Delete(n)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "venue successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
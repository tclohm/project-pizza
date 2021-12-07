package main

import (
	"fmt"
	"net/http"

	"github.com/tclohm/project-pizza/internal/data"
	"github.com/tclohm/project-pizza/internal/validator"

	_"github.com/gorilla/mux"
)

func (app *application) createVenuePizzaHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		VenueId int64 `json:"venue_id"`
		PizzaId int64 `json:"pizza_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	venuepizza := &data.VenuePizza{
		VenueId: input.venue_id,
		PizzaId: input.pizza_id,
	}

	v := validator.New()

	if data.ValidateVenuePizza(v, venuepizza); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.VenuePizzas.Insert(venuepizza)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprinf("/v1/venuepizza/%d", venuepizza.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"venuepizza": venuepizza}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}


func (app *application) updateVenuePizzaHandler(w http.ResponseWriter, r *http.Request) {
	var mux.Vars(r)
	id := vars["id"]

	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	venuepizza, err := app.models.VenuePizzas.Get(n)
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
		VenueId int64 `json:"venue_id"`
		PizzaId int64 `json:"pizza_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.VenueId != nil {
		venuepizza.VenueId = *input.VenueId
	}

	if input.PizzaId != nil {
		venuepizza.PizzaId = *input.PizzaId
	}

	v := validator.New()

	if data.ValidateVenuePizza(v, venuepizza); !v.Venue() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.VenuePizzas.Update(venue)
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

func (app *application) deleteVenuePizzaHandler(w http.ResponseWriter, r *http.Request) {
	var mux.Vars(r)
	id := vars["id"]

	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.VenuePizzas.Delete(n)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "pizza venue connection successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
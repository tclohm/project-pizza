package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tclohm/project-pizza/internal/data"
	"github.com/tclohm/project-pizza/internal/validator"

	"github.com/gorilla/mux"
)

func (app *application) createPizzaHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name 				string `json:"name"`
		Style 				string `json:"style"`
		Description 		string `json:"description"`
		Cheesiness 			float32 `json:"cheesiness"`
		Flavor 				float32 `json:"flavor"`
		Sauciness 			float32 `json:"sauciness"`
		Saltiness 			float32 `json:"saltiness"`
		Charness 			float32 `json:"charness"`
		ImageFilename 		string `json:"filename"`
		ImageContentType 	string `json:"content_type"`
		ImageLocation 		string `json:"location"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	pizza := &data.Pizza{
		Name: 				input.Name,
		Style: 				input.Style,
		Description: 		input.Description,
		Cheesiness: 		input.Cheesiness,
		Flavor: 			input.Flavor,
		Sauciness: 			input.Sauciness,
		Saltiness: 			input.Saltiness,
		Charness: 			input.Charness,
		ImageFilename: 		input.ImageFilename,
		ImageContentType: 	input.ImageContentType,
		ImageLocation: 		input.ImageLocation,
	}

	v := validator.New()

	if data.ValidatePizza(v, pizza); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showPizzaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	pizza := data.Pizza{
		ID: n,
		Name: "Mama's Fresh To Death Za",
		Style: "California-style",
		Description: "Fresh Veggies and plenty of cheese",
		Cheesiness: 3.2,
		Flavor: 5.0,
		Sauciness: 3.5,
		Saltiness: 3.0,
		Charness: 3.0,
		ImageFilename: "throwaway.jpg",
		ImageContentType: "jpg/png",
		ImageLocation: "/Users/taylor/Desktop/web/web-projects/pizza-hunter/api/uploads/throwaway.jpg",
		CreatedAt: time.Now(),
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"pizza": pizza}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
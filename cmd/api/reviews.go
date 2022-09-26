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

func (app *application) createReviewHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Style 				string 	`json:"style"`
		Price 				float32 `json:"price"`
		Description 		string 	`json:"description"`
		Cheesiness 			float32 `json:"cheesiness"`
		Flavor 				float32 `json:"flavor"`
		Sauciness 			float32 `json:"sauciness"`
		Saltiness 			float32 `json:"saltiness"`
		Charness 			float32 `json:"charness"`
		ImageId				int64 	`json:"image_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	review := &data.Review{
		Style: 				input.Style,
		Price: 				input.Price,
		Description: 		input.Description,
		Cheesiness: 		input.Cheesiness,
		Flavor: 			input.Flavor,
		Sauciness: 			input.Sauciness,
		Saltiness: 			input.Saltiness,
		Charness: 			input.Charness,
		ImageId:			input.ImageId,
	}

	v := validator.New()

	if data.ValidateReview(v, review); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Reviews.Insert(review)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/Reviews/%d", review.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"review": review}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	review, err := app.models.Reviews.Get(n)

	if err != nil {
		switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
		}
		return
	}
	

	err = app.writeJSON(w, http.StatusOK, envelope{"review": review}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	review, err := app.models.Reviews.Get(n)
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
		Style 				*string 	`json:"style"`
		Price  				*float32 	`json:"price"`
		Description 		*string 	`json:"description"`
		Cheesiness 			*float32 	`json:"cheesiness"`
		Flavor 				*float32 	`json:"flavor"`
		Sauciness 			*float32 	`json:"sauciness"`
		Saltiness 			*float32 	`json:"saltiness"`
		Charness 			*float32 	`json:"charness"`
		ImageId				*int64 	 	`json:image_id`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Style != nil {
		review.Style = *input.Style
	}

	if input.Price != nil {
		review.Price = *input.Price
	}

	if input.Description != nil {
		review.Description = *input.Description
	}

	if input.Cheesiness != nil {
		review.Cheesiness = *input.Cheesiness
	} 

	if input.Flavor != nil {
		review.Flavor = *input.Flavor
	} 	

	if input.Sauciness != nil {
		review.Sauciness = *input.Sauciness
	}

	if input.Saltiness != nil {
		review.Saltiness = *input.Saltiness
	}

	if input.Charness != nil {
		review.Charness = *input.Charness
	} 

	if input.ImageId != nil {
		review.ImageId = *input.ImageId
	}	


	v := validator.New()

	if data.ValidateReview(v, review); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Reviews.Update(review)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"review": review}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}	
}

func (app *application) deleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Reviews.Delete(n)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "Review successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listReviewsHandler(w http.ResponseWriter, r *http.Request) {

	Reviews, err := app.models.Reviews.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"Reviews": Reviews}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
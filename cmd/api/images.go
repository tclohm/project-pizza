package main

import (
	"fmt"
	"net/http"
	"strconv"
	"errors"
	"os"
	"io/ioutil"

	"github.com/tclohm/project-pizza/internal/data"
	"github.com/tclohm/project-pizza/internal/validator"

	"github.com/gorilla/mux"
)

func (app *application) createImageHandler(w http.ResponseWriter, r *http.Request) {

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	defer file.Close()

	tmpFile, err := ioutil.TempFile("uploads", "upload-*.jpg")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	defer tmpFile.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	tmpFile.Write(fileBytes)

	image := &data.Image{
		Filename: handler.Filename,
		ContentType: handler.Header["Content-Type"][0],
		Location: tmpFile.Name(),
	}

	v := validator.New()

	if data.ValidateImage(v, image); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Images.Insert(image)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/images/%d", image.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"image": image}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showImageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	image, err := app.models.Images.Get(n)

	if err != nil {
		switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
		}
		return
	}

	http.ServeFile(w, r, os.Getenv("FILEPATH") + image.Location)
}
package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tclohm/project-pizza/internal/data"

	"github.com/gorilla/mux"
)

func (app *application) createPizzaHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a pizza")
}

func (app *application) showPizzaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encounted a problem and could not process the request", http.StatusInternalServerError)
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
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
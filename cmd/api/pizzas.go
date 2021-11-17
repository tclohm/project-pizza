package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (app *application) createPizzaHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a pizza")
}

func (app *application) showPizzaHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Fprintf(w, "show the details of pizza %s\n", id)
}
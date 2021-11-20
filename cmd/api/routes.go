package main

import (
	"github.com/gorilla/mux"
)

func (app *application) routes() *mux.Router {
	// init new 
	router := mux.NewRouter()
	sub := router.PathPrefix("/v1").Subrouter()
	sub.HandleFunc("/healthcheck", app.healthcheckHandler).Methods("GET")
	sub.HandleFunc("/pizzas", app.createPizzaHandler).Methods("POST")
	sub.HandleFunc("/pizzas/{id:[0-9]+}", app.showPizzaHandler).Methods("GET")
	sub.HandleFunc("/pizzas/{id:[0-9]+}", app.updatePizzaHandler).Methods("PUT")

	return router
}
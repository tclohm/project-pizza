package main

import (
	"net/http"
	"github.com/gorilla/mux"
)

func (app *application) routes() http.Handler {
	// init new 
	router := mux.NewRouter()
	sub := router.PathPrefix("/v1").Subrouter()
	sub.HandleFunc("/healthcheck", app.healthcheckHandler).Methods("GET")
	sub.HandleFunc("/images", app.createImageHandler).Methods("POST")
	sub.HandleFunc("/images/{id:[0-9]+}", app.showImageHandler).Methods("GET")
	sub.HandleFunc("/venues", app.createVenueHandler).Methods("POST")
	sub.HandleFunc("/pizzas", app.listPizzasHandler).Methods("GET")
	sub.HandleFunc("/pizzas", app.createPizzaHandler).Methods("POST")
	sub.HandleFunc("/pizzas/{id:[0-9]+}", app.showPizzaHandler).Methods("GET")
	sub.HandleFunc("/pizzas/{id:[0-9]+}", app.updatePizzaHandler).Methods("PATCH")
	sub.HandleFunc("/pizzas/{id:[0-9]+}", app.deletePizzaHandler).Methods("DELETE")
	sub.HandleFunc("/venuepizza", app.createVenuePizzaHandler).Methods("POST")

	return app.recoverPanic(app.enableCORS(app.rateLimit(router)))
}
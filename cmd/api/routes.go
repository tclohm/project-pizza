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
	sub.HandleFunc("/venues/{id:[0-9]+}", app.showVenueHandler).Methods("GET")
	sub.HandleFunc("/reviews", app.createReviewHandler).Methods("POST")
	sub.HandleFunc("/reviews", app.listReviewsHandler).Methods("GET")
	sub.HandleFunc("/reviews/from={start}-to={end}", app.showReviewHandler).Methods("GET")	
	sub.HandleFunc("/pizzas", app.listPizzasHandler).Methods("GET")
	sub.HandleFunc("/pizzas", app.createPizzaHandler).Methods("POST")
	sub.HandleFunc("/pizzas/{id:[0-9]+}", app.showPizzaHandler).Methods("GET")
	sub.HandleFunc("/pizzas/{id:[0-9]+}", app.updatePizzaHandler).Methods("PATCH")
	sub.HandleFunc("/pizzas/{id:[0-9]+}", app.deletePizzaHandler).Methods("DELETE")
	sub.HandleFunc("/venuepizza", app.createVenuePizzaHandler).Methods("POST")
	sub.HandleFunc("/venuepizzas", app.listVenuePizzaHandler).Methods("GET")
	sub.HandleFunc("/venuepizzas/{pizzaId:[0-9]+}", app.showVenuePizzaHandler).Methods("GET")

	return app.recoverPanic(app.enableCORS(app.rateLimit(router)))
}
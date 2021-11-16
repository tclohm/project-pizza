package main

import (
	"github.com/gorilla/mux"
)

func (app *application) routes() {
	router := mux.NewRouter()

	router.HandleFunc("/healthcheck", app.healthcheckHandler)

	router.HandleFunc("/upload/image", app.PostImageHandler)
	router.HandleFunc("/post/venue", app.PostVenueHandler)

	router.HandleFunc("/image/{id}", app.GetImageHandler)

	router.HandleFunc("/post/pizza", app.PostPizzaHandler)
	router.HandleFunc("/post/venuepizza", app.PostVenuePizzaHandler)

	router.HandleFunc("/get/mypizzas", app.GetMyPizzasHandler)

	return router
}
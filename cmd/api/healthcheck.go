package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// create a map
	data := map[string]string{
		"status": "available",
		"environment": app.config.env,
		"version": version,
	}
	
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encounted a problem and could not process the request", http.StatusInternalServerError)
	}
}
package main

import (
	"log"
	"fmt"

	"github.com/tclohm/project-pizza/models"
	"github.com/tclohm/project-pizza/handlers"

)


func main() {
	db, err := models.ConnectDB()
	if err != nil {
		log.Println(err)
	}

	db.AutoMigrate(&models.Image{}, &models.Taste{}, &models.Pizza{}, &models.Venue{}, &models.VenuePizza{})

	fmt.Println(db.Migrator().CurrentDatabase())

	dbclient := handlers.DBClient{Db: db}
	
	fmt.Println("server running on 8000...")
	handlers.HandleRequests(dbclient)
}
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

	db.AutoMigrate(&models.Image{})
	db.AutoMigrate(&models.Taste{})
	db.AutoMigrate(&models.Pizza{})

	dbclient := handlers.DBClient{Db: db}
	
	fmt.Println("server running on 8000...")
	handlers.HandleRequests(dbclient)
}
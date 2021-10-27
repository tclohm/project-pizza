package models

import (
	"time"
	"fmt"
	"log"
	"os"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
)

// Time string equivalent to Date.now().toISOString in js
func GenerateISOString() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.999Z07:00")
}

func InitDB() (*sql.DB, error) {
	err := godotenv.Load()
	pg_connection_string := fmt.Sprintf("port=%s host=%s user=%s "+ 
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("PSQL_PORT"), os.Getenv("HOSTNAME"), os.Getenv("PSQL_USER"), "", os.Getenv("PSQL_DATABASE"))


	db, err := sql.Open("postgres", pg_connection_string)
	if err != nil {
		return nil, err
	}

	log.Println("Database connnected initialized.")

	return db, nil
}
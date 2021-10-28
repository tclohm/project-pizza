package models

import (
	"time"
	"fmt"
	"log"
	"os"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Time string equivalent to Date.now().toISOString in js
func GenerateISOString() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.999Z07:00")
}

func ConnectDB() (*gorm.DB, error) {
	err := godotenv.Load()
	pg_connection_string := fmt.Sprintf("host=%s user=%s password=%s  "+ 
		"dbname=%s port=%s sslmode=disable",
		os.Getenv("HOSTNAME"), os.Getenv("PSQL_USER"), "", os.Getenv("PSQL_DATABASE"), os.Getenv("PSQL_PORT"))

	sqlDB, err := sql.Open("postgres", pg_connection_string)
	
	if err != nil {
		return nil, err
	}

	log.Println("connection:", pg_connection_string)

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	log.Println("Database connnected initialized.")

	return db, nil
}
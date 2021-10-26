package models

import (
	"time"
	"fmt"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/google/uuid"
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

	stmt, err := db.Prepare(
		"CREATE TABLE IF NOT EXISTS tastes (
			taste_id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			cheesiness INTEGER NOT NULL,
			flavor INTEGER NOT NULL,
			sauciness INTEGER NOT NULL,
			saltiness INTEGER NOT NULL,
			charness INTEGER NOT NULL
			created_at TIMESTAMP NOT NULL
		)

		CREATE TABLE IF NOT EXISTS pizzas (
			pizza_id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			name VARCHAR ( 50 ) NOT NULL,
			category VARCHAR ( 50 ) NOT NULL,
			description TEXT NOT NULL,
			taste_id INT REFERENCES taste(taste_id) ON DELETE CASCADE,
			updated_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL
		)

		CREATE TABLE IF NOT EXISTS venues (
			venue_id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			name VARCHAR ( 50 ) NOT NULL,
			lat DOUBLE PRECISION,
			lon DOUBLE PRECISION,
			address VARCHAR ( 80 )
		)

		CREATE TABLE IF NOT EXISTS venuePizza (
			venue_pizza_id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
			pizza_id INT REFERENCES pizzas(pizza_id) ON DELETE CASCADE,
			venue_id INT REFERENCES venue(venue_id) ON DELETE CASCADE
		);"
	)

	if err != nil {
		return nil, err
	}

	_, err := stmt.Exec()

	if err != nil {
		return nil, err
	}

	return db, nil
}
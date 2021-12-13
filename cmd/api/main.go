package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	"context"
	"database/sql"
	"strings"

	"github.com/tclohm/project-pizza/internal/data"
	"github.com/tclohm/project-pizza/internal/jsonlog"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env string
	db struct {
		dataSource string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime string
	}
	// requests-per-second, burst values, enable/disable rate limiting
	limiter struct {
		rps 	float64
		burst 	int
		enabled bool
	}
	cors struct {
		trustedOrigins []string
	}
}

type application struct {
	config config
	logger *jsonlog.Logger
	models data.Models
}

func main() {

	// stdout streams
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	err := godotenv.Load(".env")
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	var cfg config

	// values of the env command-line flags into the config struct
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	connectionString := "host=%s user=%s dbname=%s sslmode=%s"
	connectionString = fmt.Sprintf(connectionString, os.Getenv("HOSTNAME"), os.Getenv("PSQL_USER"), os.Getenv("PSQL_DATABASE"), "disable")

	flag.StringVar(&cfg.db.dataSource, "db-datasource", connectionString, "PostgreSQL DSN")

	// limit of open connections
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	// number of idle connections in pool
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	// maximum length of time that a connection can be idle before marked as expired
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 10, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	defer db.Close()

	logger.PrintInfo("database connection pool established", nil)

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	// start server
	err = app.serve()

	if err != nil {
		logger.PrintFatal(err, nil)
	}
}


func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dataSource)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	// context with 5 second timout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// establish connection passing our context
	// if connection not established in 5 seconds
	// return error
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
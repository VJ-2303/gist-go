package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "devolopment", "Environment (devolopment | staging | production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://gistclone:your-password@localhost:5432/gistclone?sslmode=disable", "PostgreSQL DSN")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: cfg,
		logger: logger,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
	}

	logger.Printf("starting %s server on %d", app.config.env, app.config.port)
	err := srv.ListenAndServe()
	logger.Fatal(err)

	fmt.Println("Starting the gistclone API...")
}

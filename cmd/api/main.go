package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config

	cfg.port = 4000
	cfg.env = "devolopment"

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

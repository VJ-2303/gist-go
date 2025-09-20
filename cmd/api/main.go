package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/vj-2303/gist-go/internal/data"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
	jwt struct {
		secret string
	}
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "devolopment", "Environment (devolopment | staging | production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://gistclone:your-password@localhost:5432/gistclone?sslmode=disable", "PostgreSQL DSN")

	flag.StringVar(&cfg.jwt.secret, "jwt-secret", "", "JWT Secret String")

	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Println("database connection pool established")

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
	}

	logger.Printf("starting %s server on %d", app.config.env, app.config.port)
	err = srv.ListenAndServe()
	logger.Fatal(err)

	fmt.Println("Starting the gistclone API...")
}

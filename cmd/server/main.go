package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/danielrpof/go-automation-runner/internal/api"
	"github.com/danielrpof/go-automation-runner/internal/auth"
	"github.com/danielrpof/go-automation-runner/internal/store"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY environment variable is not set")
	}

	db, err := sql.Open("sqlite3", "./jobs.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	schema := `
	CREATE TABLE IF NOT EXISTS jobs (
		id TEXT PRIMARY KEY,
		command TEXT NOT NULL,
		status TEXT NOT NULL,
		stdout TEXT,
		stderr TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);
	`

	if _, err := db.Exec(schema); err != nil {
		log.Fatal(err)
	}

	store := store.NewSQLiteStore(db)
	mux := http.NewServeMux()

	//public
	mux.HandleFunc("/health", api.HealthHandler)

	//protected
	protected := auth.APIKeyMiddleware(apiKey)
	mux.Handle("/jobs", protected(api.JobsHandler(store)))
	mux.Handle("/jobs/", protected(api.JobByIDHandler(store)))

	addr := ":8080"
	log.Println("Starting server on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

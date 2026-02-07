package main

import (
	"log"
	"net/http"
	"os"

	"github.com/danielrpof/go-automation-runner/internal/api"
	"github.com/danielrpof/go-automation-runner/internal/auth"
	"github.com/danielrpof/go-automation-runner/internal/job"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY environment variable is not set")
	}

	jobStore := job.NewStore()

	mux := http.NewServeMux()

	//public
	mux.HandleFunc("/health", api.HealthHandler)

	//protected
	protected := auth.APIKeyMiddleware(apiKey)
	mux.Handle("/jobs", protected(api.JobsHandler(jobStore)))
	mux.Handle("/jobs/", protected(api.JobByIDHandler(jobStore)))

	addr := ":8080"
	log.Println("Starting server on", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

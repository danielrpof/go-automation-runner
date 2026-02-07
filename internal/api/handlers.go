package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/danielrpof/go-automation-runner/internal/job"
	"github.com/google/uuid"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func JobsHandler(store *job.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CreateJobRequest

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if req.Command == "" {
			http.Error(w, "command is required", http.StatusBadRequest)
			return
		}

		newJob := &job.Job{
			ID:        uuid.New().String(),
			Command:   req.Command,
			Status:    job.StatusPending,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		store.Add(newJob)
		go job.Run(newJob)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newJob)
	}

}

func JobByIDHandler(store *job.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Expected path: /jobs/{id}
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) != 3 {
			http.Error(w, "invalid job id", http.StatusBadRequest)
			return
		}

		jobID := parts[2]
		j, ok := store.Get(jobID)
		if !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(j)
	}
}

type CreateJobRequest struct {
	Command string `json:"command"`
}

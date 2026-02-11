package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/danielrpof/go-automation-runner/internal/job"
	"github.com/danielrpof/go-automation-runner/internal/store"
	"github.com/google/uuid"
)

// AllowedJobs defines the allowlist of commands that can be executed
var AllowedJobs = map[string]string{
	"disk-check":   "df -h",
	"memory-check": "free -m",
	"uptime":       "uptime",
	"list-procs":   "ps aux | head -20",
	"say-hello":    "echo Hello from automation runner!",
	"say_hello":    "echo Hello from automation runner!",
	"echo-test":    "echo This is a test",
	"date":         "date",
	"sleep":        "powershell -Command Start-Sleep -Seconds 40",
	"timeout":      "powershell -Command Start-Sleep -Seconds 40",
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func JobsHandler(store store.JobStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodPost:
			var req CreateJobRequest

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}

			if req.Job == "" {
				http.Error(w, "job is required", http.StatusBadRequest)
				return
			}

			cmdArgs, ok := AllowedJobs[req.Job]
			if !ok {
				http.Error(w, "job not allowed", http.StatusBadRequest)
				return
			}

			newJob := &job.Job{
				ID:        uuid.New().String(),
				Command:   cmdArgs,
				Status:    job.StatusPending,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			if err := store.Add(newJob); err != nil {
				http.Error(w, "failed to create job", http.StatusInternalServerError)
				return
			}
			go job.Run(newJob, store)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(newJob)

		case http.MethodGet:
			jobs, err := store.List()
			if err != nil {
				http.Error(w, "failed to list jobs", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(jobs)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func JobByIDHandler(store store.JobStore) http.HandlerFunc {
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
		j, err := store.Get(jobID)
		if err != nil {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(j)
	}
}

type CreateJobRequest struct {
	Job string `json:"job"`
}

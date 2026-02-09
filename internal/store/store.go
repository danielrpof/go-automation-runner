package store

import "github.com/danielrpof/go-automation-runner/internal/job"

type JobStore interface {
	Add(job *job.Job) error
	Get(id string) (*job.Job, error)
	List() ([]*job.Job, error)
}

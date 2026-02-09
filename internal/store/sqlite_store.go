package store

import (
	"database/sql"
	"time"

	"github.com/danielrpof/go-automation-runner/internal/job"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(db *sql.DB) *SQLiteStore {
	return &SQLiteStore{db: db}
}

// Add creates a new job in the database
func (s *SQLiteStore) Add(j *job.Job) error {
	_, err := s.db.Exec(`
		INSERT INTO jobs (id, command, status, stdout, stderr, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		j.ID,
		j.Command,
		j.Status,
		j.Stdout,
		j.Stderr,
		j.CreatedAt,
		j.UpdatedAt,
	)

	return err
}

func (s *SQLiteStore) Create(j *job.Job) error {
	_, err := s.db.Exec(`
		INSERT INTO jobs (id, command, status, stdout, stderr, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`,
		j.ID,
		j.Command,
		j.Status,
		j.Stdout,
		j.Stderr,
		j.CreatedAt,
		j.UpdatedAt,
	)

	return err
}

func (s *SQLiteStore) Update(j *job.Job) error {
	_, err := s.db.Exec(`
		UPDATE jobs
		SET status = ?, stdout = ?, stderr = ?, updated_at = ?
		WHERE id = ?
	`,
		j.Status,
		j.Stdout,
		j.Stderr,
		time.Now(),
		j.ID,
	)

	return err
}

func (s *SQLiteStore) Get(id string) (*job.Job, error) {
	row := s.db.QueryRow(`
		SELECT id, command, status, stdout, stderr, created_at, updated_at
		FROM jobs WHERE id = ?
	`, id)

	var j job.Job
	err := row.Scan(
		&j.ID,
		&j.Command,
		&j.Status,
		&j.Stdout,
		&j.Stderr,
		&j.CreatedAt,
		&j.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &j, nil
}

func (s *SQLiteStore) List() ([]*job.Job, error) {
	rows, err := s.db.Query(`
		SELECT id, command, status, stdout, stderr, created_at, updated_at
		FROM jobs
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []*job.Job
	for rows.Next() {
		var j job.Job
		err := rows.Scan(
			&j.ID,
			&j.Command,
			&j.Status,
			&j.Stdout,
			&j.Stderr,
			&j.CreatedAt,
			&j.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, &j)
	}

	return jobs, nil
}

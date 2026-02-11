package job

import (
	"bytes"
	"context"
	"os/exec"
	"runtime"
	"time"
)

type JobUpdater interface {
	Update(job *Job) error
}

func Run(job *Job, store JobUpdater) {
	job.Status = StatusRunning
	job.UpdatedAt = time.Now()
	store.Update(job)

	var cmd *exec.Cmd
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", job.Command)
	} else {
		cmd = exec.CommandContext(ctx, "sh", "-c", job.Command)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	job.Stdout = stdout.String()
	job.Stderr = stderr.String()

	if err != nil {
		job.Status = StatusFailed
		// Check if timeout occurred
		if ctx.Err() == context.DeadlineExceeded {
			job.Stdout = "Command execution timed out after 30 seconds"
		}
	} else {
		job.Status = StatusCompleted
	}

	job.UpdatedAt = time.Now()
	store.Update(job)
}

package job

import (
	"bytes"
	"os/exec"
	"runtime"
	"time"
)

func Run(job *Job) {
	job.Status = StatusRunning
	job.UpdatedAt = time.Now()

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", job.Command)
	} else {
		cmd = exec.Command("sh", "-c", job.Command)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	job.Stdout = stdout.String()
	job.Stderr = stderr.String()

	if err != nil {
		job.Status = StatusFailed
	} else {
		job.Status = StatusCompleted
	}

	job.UpdatedAt = time.Now()
}

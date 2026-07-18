package application

import (
	"context"
	"os/exec"
	"regexp"
	"time"
)

type TestRunResult struct {
	Success bool   `json:"success"`
	Passed  int    `json:"passed"`
	Failed  int    `json:"failed"`
	Output  string `json:"output"`
}

var passLine = regexp.MustCompile(`(?m)^--- PASS:`)
var failLine = regexp.MustCompile(`(?m)^--- FAIL:`)

// RunBackendTests shells out to `go test ./... -v` with a bounded timeout.
// This is a single-operator internal dev tool, not a job queue: it blocks the
// HTTP request for the duration of the run.
type RunBackendTests struct {
	WorkDir string
}

func (uc *RunBackendTests) Execute(ctx context.Context) (TestRunResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "test", "./...", "-v")
	cmd.Dir = uc.WorkDir
	output, runErr := cmd.CombinedOutput()

	result := TestRunResult{
		Success: runErr == nil,
		Passed:  len(passLine.FindAllIndex(output, -1)),
		Failed:  len(failLine.FindAllIndex(output, -1)),
		Output:  string(output),
	}
	return result, nil
}

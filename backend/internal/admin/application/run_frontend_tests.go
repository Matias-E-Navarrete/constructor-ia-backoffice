package application

import (
	"context"
	"os/exec"
	"time"
)

// RunFrontendTests shells out to `npx playwright test` with a bounded
// timeout. Playwright's own summary line ("5 passed", "2 failed, 3 passed")
// is left in Output for the admin UI to display verbatim rather than
// re-parsing its reporter format here.
type RunFrontendTests struct {
	WorkDir string
}

func (uc *RunFrontendTests) Execute(ctx context.Context) (TestRunResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "npx", "playwright", "test")
	cmd.Dir = uc.WorkDir
	output, runErr := cmd.CombinedOutput()

	return TestRunResult{
		Success: runErr == nil,
		Output:  string(output),
	}, nil
}

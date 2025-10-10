package shell

import (
	"bytes"
	"io"
	"os"
	"context"
	"os/exec"
	"regexp"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy/internal/application/ports"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

type Executor struct{}

func NewExecutor() ports.CommandExecutor {
	return &Executor{}
}

func (e *Executor) CreateWorkDir(workdirs ...string) string {
	return filepath.Join(workdirs...)
}

func (e *Executor) Execute(ctx context.Context, workdir, command string) (log string, exitCode int, err error) {
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	if workdir != "" {
		cmd.Dir = workdir
	}

	var out bytes.Buffer
	multiOutput := io.MultiWriter(os.Stdout, &out)

	cmd.Stdout = multiOutput
	cmd.Stderr = multiOutput

	runErr := cmd.Run()
	log = out.String()

	if runErr != nil {
		if exitErr, ok := runErr.(*exec.ExitError); ok {
			return log, exitErr.ExitCode(), nil
		}
		return log, -1, runErr
	}

	log = ansiRegex.ReplaceAllString(log, "")
	return log, 0, nil
}

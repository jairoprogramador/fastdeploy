package application

import (
	"bytes"
	//"io"
	//"os"
	"context"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/jairoprogramador/fastdeploy-core/internal/application/ports"
)

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

type ExecCommandService struct{}

func NewExecCommandService() ports.CommandService {
	return &ExecCommandService{}
}

func (e *ExecCommandService) CreateWorkDir(workdirs ...string) string {
	return filepath.Join(workdirs...)
}

func (e *ExecCommandService) Run(ctx context.Context, workdir, command string) (string, int, error) {
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	if workdir != "" {
		cmd.Dir = workdir
	}

	var buffer bytes.Buffer
	//multiOutput := io.MultiWriter(os.Stdout, &out)

	//cmd.Stdout = multiOutput
	//cmd.Stderr = multiOutput

	cmd.Stdout = &buffer
	cmd.Stderr = &buffer

	err := cmd.Run()
	output := buffer.String()

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return output, exitErr.ExitCode(), nil
		}
		return output, -1, err
	}

	output = ansiRegex.ReplaceAllString(output, "")
	return output, 0, nil
}

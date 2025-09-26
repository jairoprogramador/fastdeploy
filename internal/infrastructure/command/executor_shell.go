package command

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/jairoprogramador/fastdeploy/internal/domain/command/port"
)

type ExecutorShell struct{}

func NewExecutorShell() port.ExecutorPort {
	return &ExecutorShell{}
}

func (e *ExecutorShell) Run(command string, workdir string) (string, error) {
	commandExec := exec.Command("sh", "-c", command)

	var output bytes.Buffer
	multiOutput := io.MultiWriter(os.Stdout, &output)

	commandExec.Dir = workdir
	fmt.Printf("   -> Command exec dir: %s\n", workdir)
	commandExec.Stdout = multiOutput
	commandExec.Stderr = multiOutput

	if err := commandExec.Run(); err != nil {
		return "", err
	}

	return output.String(), nil
}

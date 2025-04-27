package tools

import (
	"fmt"
	"bytes"
	"os/exec"
	"strings"
	"deploy/internal/interface/presenter"
)

func ExecuteCommand(command string, args ...string) (string, error) {
	output := fmt.Sprintf("command executed: '%s %s'", command, strings.Join(args, " "))
	fmt.Println(output)

	var stdoutBuf, stderrBuf bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Stdout = &stdoutBuf
    cmd.Stderr = &stderrBuf

	done := make(chan bool)
	go presenter.ShowLoader(done)

	err := cmd.Run()

	done <- true
    fmt.Println()

	if err != nil {
		if stdoutBuf.Len() > 0 {
			fmt.Println("Standard Output:")
			fmt.Println(stdoutBuf.String())
		}
		if stderrBuf.Len() > 0 {
			fmt.Println("Standard Error:")
			fmt.Println(stderrBuf.String())
		}
		stdoutBuf.Reset()
    	stderrBuf.Reset()

		fmt.Println(output)
		return "", err
	}

	stderrBuf.Reset()
	return stdoutBuf.String(), nil
}
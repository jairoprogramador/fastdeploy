package tools

import (
	"context"
	"fmt"
	"bytes"
	"os/exec"
	"strings"
	"time"
)

func ExecuteCommand(command string, args ...string) (string, error) {
	output := fmt.Sprintf("command executed: '%s %s'", command, strings.Join(args, " "))
	fmt.Println(output)

	var stdoutBuf, stderrBuf bytes.Buffer

	cmd := exec.Command(command, args...)
	cmd.Stdout = &stdoutBuf
    cmd.Stderr = &stderrBuf

	//done := make(chan bool)
	//go presenter.ShowLoader(done)

	err := cmd.Run()

	//done <- true
    //fmt.Println()

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

func ExecuteCommandWithContext(ctx context.Context, command string, args ...string) (string, error) {
	output := fmt.Sprintf("command executed with context: '%s %s'", command, strings.Join(args, " "))
	fmt.Println(output)

	var stdoutBuf, stderrBuf bytes.Buffer

	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	var cancel context.CancelFunc
	_, hasDeadline := ctx.Deadline()
	if !hasDeadline {
		ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
	}

	err := cmd.Run()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout al ejecutar comando: %s", output)
		}
		
		if ctx.Err() == context.Canceled {
			return "", fmt.Errorf("comando cancelado: %s", output)
		}
		
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
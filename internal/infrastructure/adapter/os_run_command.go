package adapter

import (
	"bytes"
	"context"
	"deploy/internal/domain/model"
	"deploy/internal/domain/port"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Constants for configuration
const (
	DefaultCommandTimeout = 10 * time.Minute
)

// osRunCommand implements the RunCommand interface
type osRunCommand struct{}

// NewOsRunCommand creates a new instance of the command executor
func NewOsRunCommand() port.RunCommand {
	return &osRunCommand{}
}

// Run executes a command and returns its result
func (r *osRunCommand) Run(ctx context.Context, cmdExec string) model.InfrastructureResponse {
	if cmdExec == "" {
		return model.NewErrorResponse(fmt.Errorf("the command is empty"))
	}

	ctx = ensureContextTimeout(ctx)
	command, args := parseCommand(cmdExec)

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := prepareCommand(ctx, command, args, &stdoutBuf, &stderrBuf)

	cmdDetails := formatCommandDetails(command, args)

	err := cmd.Run()
	if err != nil {
		return handleCommandError(ctx, err, cmdDetails, &stdoutBuf, &stderrBuf)
	}

	return model.NewResponseWithDetails(
		strings.TrimSpace(stdoutBuf.String()), 
		cmdDetails,
	)
}

// ensureContextTimeout ensures the context has a timeout
func ensureContextTimeout(ctx context.Context) context.Context {
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		newCtx, _ := context.WithTimeout(ctx, DefaultCommandTimeout)
		return newCtx
	}
	return ctx
}

// parseCommand splits a command string into command and arguments
func parseCommand(cmdExec string) (string, []string) {
	parts := strings.Fields(cmdExec)
	command := parts[0]
	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}
	return command, args
}

// prepareCommand creates and configures an exec.Cmd instance
func prepareCommand(ctx context.Context, command string, args []string, stdout, stderr *bytes.Buffer) *exec.Cmd {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd
}

// formatCommandDetails creates a string with command execution details
func formatCommandDetails(command string, args []string) string {
	return fmt.Sprintf("Running command: '%s %s'\n", command, strings.Join(args, " "))
}

// handleCommandError processes command execution errors
func handleCommandError(ctx context.Context, err error, cmdDetails string, stdout, stderr *bytes.Buffer) model.InfrastructureResponse {
	var message strings.Builder
	message.WriteString(cmdDetails)
	message.WriteString(fmt.Sprintf("Error running command: %v\n", err))

	// Check for context errors first
	if ctxErr := ctx.Err(); ctxErr != nil {
		if ctxErr == context.DeadlineExceeded {
			message.WriteString(fmt.Sprintf("Timeout (context deadline exceeded): %v\n", ctxErr))
			return model.NewErrorResponseWithDetails(ctxErr, message.String())
		}

		if ctxErr == context.Canceled {
			message.WriteString(fmt.Sprintf("Command canceled: %v\n", ctxErr))
			return model.NewErrorResponseWithDetails(ctxErr, message.String())
		}
	}

	// Add output information if available
	appendOutputToMessage(&message, stdout, stderr)

	return model.NewErrorResponseWithDetails(err, message.String())
}

// appendOutputToMessage adds command output to the error message
func appendOutputToMessage(message *strings.Builder, stdout, stderr *bytes.Buffer) {
	stdOutput := strings.TrimSpace(stdout.String())
	if stdOutput != "" {
		message.WriteString(fmt.Sprintf("Standard Output:\n%s\n", stdOutput))
	}

	stdError := strings.TrimSpace(stderr.String())
	if stdError != "" {
		message.WriteString(fmt.Sprintf("Standard Error:\n%s\n", stdError))
	}
}

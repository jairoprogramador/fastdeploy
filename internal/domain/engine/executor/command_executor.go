package executor

import (
	"context"
	"deploy/internal/domain/engine/condition"
	"deploy/internal/domain/engine/model"
	"deploy/internal/domain/model/logger"
	"deploy/internal/domain/port"
	"fmt"
	"strings"
)

// Condition types that require command output
var typeConditionsEvaluate = []condition.TypeCondition{
	condition.Equals,
	condition.Contains,
	condition.Matches,
}

// CommandExecutor handles command execution and condition evaluation
type CommandExecutor struct {
	baseExecutor     *StepExecutor
	commandRunner    port.RunCommand
	conditionFactory *condition.EvaluatorFactory
	variables        *model.StoreEntity
	logger           *logger.Logger
}

// NewCommandExecutor creates a new command executor instance
func NewCommandExecutor(
	logger *logger.Logger,
	baseExecutor *StepExecutor,
	variables *model.StoreEntity,
	commandRunner port.RunCommand,
	conditionFactory *condition.EvaluatorFactory,
) Executor {
	return &CommandExecutor{
		logger:           logger,
		baseExecutor:     baseExecutor,
		variables:        variables,
		commandRunner:    commandRunner,
		conditionFactory: conditionFactory,
	}
}

// Execute runs the command defined in the step and evaluates any conditions
func (e *CommandExecutor) Execute(ctx context.Context, step model.Step) error {
	ctx, cancel := e.baseExecutor.prepareContext(ctx, step)
	defer cancel()

	return e.baseExecutor.handleRetry(step, func() error {
		e.variables.PushScope(step.Variables)
		defer e.variables.PopScope()

		commandOutput, err := e.runCommand(ctx, step)
		if err != nil {
			return err
		}

		if step.If == "" {
			return nil
		}

		return e.processCondition(step, commandOutput)
	})
}

// runCommand executes the command and returns its output
func (e *CommandExecutor) runCommand(ctx context.Context, step model.Step) (string, error) {
	response := e.commandRunner.Run(ctx, step.Command)
	if !response.IsSuccess() {
		e.logger.ErrorSystemMessage(response.Details, response.Error)
		e.logger.Error(response.Error)
		return "", response.Error
	}

	return response.Result.(string), nil
}

// processCondition validates the command output against the condition
func (e *CommandExecutor) processCondition(step model.Step, commandOutput string) error {
	conditionType := e.getConditionType(step.If)

	// Check if output is required but empty
	if e.requiresEvaluate(conditionType) && commandOutput == "" {
		return e.createEmptyOutputError(conditionType, step.Name)
	}

	return e.evaluateCondition(step.If, commandOutput, step)
}

// getConditionType extracts the condition type from the condition string
func (e *CommandExecutor) getConditionType(conditionStr string) string {
	parts := strings.SplitN(conditionStr, ":", 2)
	return parts[0]
}

// requiresEvaluate checks if the condition type requires command output
func (e *CommandExecutor) requiresEvaluate(conditionType string) bool {
	for _, requiredType := range typeConditionsEvaluate {
		if string(requiredType) == conditionType {
			return true
		}
	}
	return false
}

// createEmptyOutputError creates an error for empty command output
func (e *CommandExecutor) createEmptyOutputError(conditionType, stepName string) error {
	message := fmt.Sprintf("the command did not return a value for condition %s in step %s",
		conditionType, stepName)
	return e.logger.NewError(message)
}

func (e *CommandExecutor) evaluateCondition(conditionStr string, output string, step model.Step) error {
	evaluator := e.conditionFactory.CreateEvaluator(conditionStr)

	if isCorrect := evaluator.Evaluate(output); !isCorrect {
		message := fmt.Sprintf("the condition %s in step %s is not met", conditionStr, step.Name)
		return e.logger.NewError(message)
	}
	return nil
}

package executor

import (
	"context"
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/model"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/condition"
	"github.com/jairoprogramador/fastdeploy/internal/domain/port"
	"strings"
)

var typeConditionsEvaluate = []condition.TypeCondition{
	condition.Equals,
	condition.Contains,
	condition.Matches,
}

type CommandExecutor struct {
	baseExecutor     *BaseExecutor
	commandRunner    port.CommandPort
	conditionFactory *condition.EvaluatorFactory
}

func NewCommandExecutor(
	baseExecutor *BaseExecutor,
	commandRunner port.CommandPort,
	conditionFactory *condition.EvaluatorFactory,
) Executor {
	return &CommandExecutor{
		baseExecutor:     baseExecutor,
		commandRunner:    commandRunner,
		conditionFactory: conditionFactory,
	}
}

func (e *CommandExecutor) Execute(ctx context.Context, step model.Step) error {
	ctx, cancel := e.baseExecutor.prepareContext(ctx, step)
	defer cancel()

	return e.baseExecutor.handleRetry(step, func() error {
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

func (e *CommandExecutor) runCommand(ctx context.Context, step model.Step) (string, error) {
	response := e.commandRunner.Run(ctx, step.Command)
	if !response.IsSuccess() {
		return "", response.Error
	}

	return response.Result.(string), nil
}

func (e *CommandExecutor) processCondition(step model.Step, commandOutput string) error {
	conditionType := e.getConditionType(step.If)

	if e.requiresEvaluateValue(conditionType) && commandOutput == "" {
		return e.createEmptyOutputError(conditionType, step.Name)
	}

	return e.evaluateCondition(step.If, commandOutput, step)
}

func (e *CommandExecutor) getConditionType(conditionStr string) string {
	parts := strings.SplitN(conditionStr, ":", 2)
	return parts[0]
}

func (e *CommandExecutor) requiresEvaluateValue(conditionType string) bool {
	for _, requiredType := range typeConditionsEvaluate {
		if string(requiredType) == conditionType {
			return true
		}
	}
	return false
}

func (e *CommandExecutor) createEmptyOutputError(conditionType, stepName string) error {
	message := fmt.Sprintf("the command did not return a value for condition %s in step %s",
		conditionType, stepName)
	return fmt.Errorf(message)
}

func (e *CommandExecutor) evaluateCondition(conditionStr string, output string, step model.Step) error {
	evaluator := e.conditionFactory.CreateEvaluator(conditionStr)

	if isCorrect := evaluator.Evaluate(output); !isCorrect {
		message := fmt.Sprintf("the condition %s in step %s is not met", conditionStr, step.Name)
		return fmt.Errorf(message)
	}
	return nil
}

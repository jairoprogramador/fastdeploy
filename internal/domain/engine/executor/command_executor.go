package executor

import (
	"context"
	condition2 "deploy/internal/domain/engine/condition"
	"deploy/internal/domain/model"
	"deploy/internal/domain/model/logger"
	"deploy/internal/domain/port"
	"fmt"
	"strings"
)

type CommandExecutor struct {
	baseExecutor     *BaseExecutor
	commandRunner    port.ExecutorServiceInterface
	conditionFactory *condition2.ConditionFactory
	variables        *model.VariableStore
	logger           *logger.Logger
}

func NewCommandExecutor(
	logger *logger.Logger,
	baseExecutor *BaseExecutor,
	variables *model.VariableStore,
	commandRunner port.ExecutorServiceInterface,
	conditionFactory *condition2.ConditionFactory,
) StepExecutorInterface {
	return &CommandExecutor{
		logger:           logger,
		baseExecutor:     baseExecutor,
		variables:        variables,
		commandRunner:    commandRunner,
		conditionFactory: conditionFactory,
	}
}

func (e *CommandExecutor) Execute(ctx context.Context, step model.Step) error {
	ctx, cancel := e.baseExecutor.prepareContext(ctx, step)
	defer cancel()

	return e.baseExecutor.handleRetry(step, func() error {
		e.variables.PushScope(step.Variables)
		defer e.variables.PopScope()

		response := e.commandRunner.Run(ctx, step.Command)
		if !response.IsSuccess() {
			e.logger.ErrorSystemMessage(response.Details, response.Error)
			e.logger.Error(response.Error)
			return response.Error
		}

		output := response.Result.(string)

		if step.If != "" {
			parts := strings.SplitN(step.If, ":", 2)
			if parts[0] == string(condition2.Equals) || parts[0] == string(condition2.Contains) || parts[0] == string(condition2.Matches) {
				if output == "" {
					message := fmt.Sprintf("the command not return a value for condition %s in step %s", parts[0], step.Name)
					return e.logger.NewError(message)
				}
			}

			if err := e.evaluateCondition(step.If, output, step); err != nil {
				return err
			}
		}
		return nil
	})
}

func (e *CommandExecutor) evaluateCondition(conditionStr string, output string, step model.Step) error {
	evaluator := e.conditionFactory.CreateEvaluator(conditionStr)

	if isCorrect := evaluator.Evaluate(output); !isCorrect {
		message := fmt.Sprintf("the condition %s in step %s is not met", conditionStr, step.Name)
		return e.logger.NewError(message)
	}
	return nil
}

package executor

import (
	"context"
	"deploy/internal/domain/condition"
	"deploy/internal/domain/model"
	"deploy/internal/domain/service"
	"fmt"
)

type CommandExecutor struct {
	baseExecutor     *BaseExecutor
	commandRunner    service.ExecutorServiceInterface
	conditionFactory *condition.ConditionFactory
	variables        *model.VariableStore
}

func NewCommandExecutor(
	baseExecutor *BaseExecutor,
	variables *model.VariableStore,
	commandRunner service.ExecutorServiceInterface,
	conditionFactory *condition.ConditionFactory,
) StepExecutorInterface {
	return &CommandExecutor{
		baseExecutor:     baseExecutor,
		variables:        variables,
		commandRunner:    commandRunner,
		conditionFactory: conditionFactory,
	}
}

func (e *CommandExecutor) Execute(ctx context.Context, step model.Step) (string, error) {
	ctx, cancel := e.baseExecutor.prepareContext(ctx, step)
	defer cancel()

	return e.baseExecutor.handleRetry(step, func() (string, error) {
		e.variables.PushScope(step.Variables)
		defer e.variables.PopScope()

		output, err := e.commandRunner.Run(ctx, step.Command)
		if err != nil {
			return "", fmt.Errorf("error ejecutando comando: %v", err)
		}

		if step.If != "" {
			if err := e.evaluateCondition(step.If, output); err != nil {
				fmt.Println("error en la evaluación de la condición: ", err)
				return "", err
			}
		}
		return "", nil
	})
}

func (e *CommandExecutor) evaluateCondition(conditionStr string, output string) error {
	evaluator, err := e.conditionFactory.CreateEvaluator(conditionStr, output)
	if err != nil {
		return fmt.Errorf("tipo de condición no soportado: %v", err)
	}

	result, err := evaluator.Evaluate(output)
	if err != nil {
		return fmt.Errorf("error evaluando condición: %v", err)
	}

	if !result {
		return fmt.Errorf("la condición no se cumplió: %s", conditionStr)
	}

	return nil
}

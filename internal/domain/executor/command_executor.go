package executor

import (
	"context"
	"deploy/internal/domain/condition"
	"deploy/internal/domain/model"
	"deploy/internal/domain/variable"
	"deploy/internal/domain/service"
	"fmt"
)

type CommandExecutor struct {
	baseExecutor *BaseExecutor
	commandRunner service.ExecutorServiceInterface
	conditionFactory *condition.ConditionFactory
	variables *variable.VariableStore
}

func GetCommandExecutor(variables *variable.VariableStore) *CommandExecutor {
	return &CommandExecutor {
		baseExecutor: GetBaseExecutor(),
		variables: variables,
		commandRunner: service.GetExecutorService(),
		conditionFactory: condition.GetConditionFactory(),
	}
}

func (e *CommandExecutor) Execute(ctx context.Context, step model.Step) error {
	ctx, cancel := e.baseExecutor.prepareContext(ctx, step)
	defer cancel()

	return e.baseExecutor.handleRetry(step, func() error {
		e.variables.PushScope(step.Variables)
		defer e.variables.PopScope()

		fmt.Printf("---------------%s-----------------\n", step.Name)
		output, err := e.commandRunner.Run(ctx, step.Command)
		if err != nil {
			return fmt.Errorf("error ejecutando comando: %v", err)
		}

		if step.If != "" {
			if err := e.evaluateCondition(step.If, output); err != nil {
				return err
			}
		}
		fmt.Printf("Salida del comando: \n%s", output)
		return nil
	})
}

func (e *CommandExecutor) evaluateCondition(conditionStr string, output string) error {
	fmt.Println("condition =", conditionStr)
	//fmt.Println(output)
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

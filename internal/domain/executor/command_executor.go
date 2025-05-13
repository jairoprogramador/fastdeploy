package executor

import (
	"context"
	"deploy/internal/domain/condition"
	"deploy/internal/domain/model"
	"deploy/internal/domain/variable"
	"fmt"
)

type CommandExecutor struct {
	BaseExecutor
	commandRunner    CommandRunner
	conditionFactory *condition.ConditionFactory
	variables *variable.VariableStore
}

func GetCommandExecutor(variables *variable.VariableStore) *CommandExecutor {
	return &CommandExecutor {
		BaseExecutor: BaseExecutor {},
		variables: variables,
		commandRunner:    GetCommandRunner(),
		conditionFactory: condition.GetConditionFactory(),
	}
}

func (e *CommandExecutor) Execute(ctx context.Context, step model.Step) error {
	ctx, cancel := e.prepareContext(ctx, step)
	defer cancel()

	return e.handleRetry(step, func() error {
		// Preparar variables locales
		e.variables.PushScope(step.Variables)
		defer e.variables.PopScope()

		// Ejecutar comando
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
		// Registrar la salida del comando
		fmt.Printf("Salida del comando: \n%s", output)
		return nil
	})
}

// evaluateCondition evalúa una condición sobre la salida de un comando
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

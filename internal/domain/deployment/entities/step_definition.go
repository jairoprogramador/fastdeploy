package entities

import (
	"errors"
	"fmt"
	"slices"

	shared "github.com/jairoprogramador/fastdeploy-core/internal/domain/shared"

	depVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/vos"
)

type StepDefinition struct {
	name     string
	triggers []depVos.Trigger
	commands []depVos.CommandDefinition
	variables []depVos.Variable
}

func NewStepDefinition(
	name string,
	triggers []depVos.Trigger,
	commands []depVos.CommandDefinition,
	variables []depVos.Variable) (StepDefinition, error) {

	if name == "" {
		return StepDefinition{}, errors.New("el nombre de la definición no puede estar vacío")
	}

	if len(commands) == 0 {
		return StepDefinition{}, errors.New("la definición debe tener al menos un comando")
	}

	if len(triggers) > 0 {
		validTriggers := shared.ScopesValid()
		for _, trigger := range triggers {
			if !slices.Contains(validTriggers, trigger.String()) {
				return StepDefinition{}, fmt.Errorf("trigger no válido: %s", trigger.String())
			}
		}

		triggerNames := make(map[string]struct{})
		for _, trigger := range triggers {
			if _, exists := triggerNames[trigger.String()]; exists {
				return StepDefinition{}, errors.New("trigger duplicado encontrado")
			}
			triggerNames[trigger.String()] = struct{}{}
		}
	} else {
		triggersCommon := []depVos.Trigger{
			depVos.TriggerFromString(shared.ScopeRecipe),
			depVos.TriggerFromString(shared.ScopeCode),
			depVos.TriggerFromString(shared.ScopeVars),
		}

		if name == shared.StepSupply {
			triggers = []depVos.Trigger{
				depVos.TriggerFromString(shared.ScopeRecipe),
				depVos.TriggerFromString(shared.ScopeVars),
			}
		}

		if name == shared.StepTest {
			triggers = triggersCommon
		}

		if name == shared.StepPackage {
			triggers = triggersCommon
		}
		if name == shared.StepDeploy {
			triggers = triggersCommon
		}
	}

	commandNames := make(map[string]struct{})
	for _, command := range commands {
		if _, exists := commandNames[command.Name()]; exists {
			return StepDefinition{}, fmt.Errorf("command %s duplicado", command.Name())
		}
		commandNames[command.Name()] = struct{}{}
	}

	cmds := make(map[string]struct{})
	for _, command := range commands {
		if _, exists := cmds[command.Cmd()]; exists {
			return StepDefinition{}, fmt.Errorf("comando duplicado: %s", command.Cmd())
		}
		cmds[command.Cmd()] = struct{}{}
	}

	variableNames := make(map[string]struct{})
	for _, variable := range variables {
		if _, exists := variableNames[variable.Name()]; exists {
			return StepDefinition{}, fmt.Errorf("variable duplicada: %s", variable.Name())
		}
		variableNames[variable.Name()] = struct{}{}
	}

	variablesValues := make(map[string]struct{})
	for _, variable := range variables {
		if _, exists := variablesValues[variable.Value()]; exists {
			return StepDefinition{}, fmt.Errorf("valor de variable duplicado: %s", variable.Value())
		}
		variablesValues[variable.Value()] = struct{}{}
	}

	triggersCopy := make([]depVos.Trigger, len(triggers))
	copy(triggersCopy, triggers)

	commandsCopy := make([]depVos.CommandDefinition, len(commands))
	copy(commandsCopy, commands)

	variablesCopy := make([]depVos.Variable, len(variables))
	copy(variablesCopy, variables)

	return StepDefinition{
		name:     name,
		triggers: triggersCopy,
		commands: commandsCopy,
		variables: variablesCopy,
	}, nil
}

func (sd StepDefinition) Name() string {
	return sd.name
}

func (sd StepDefinition) TriggersInt() []int {
	triggersInt := make([]int, len(sd.triggers))
	for i, trigger := range sd.triggers {
		triggersInt[i] = int(trigger)
	}
	return triggersInt
}

func (sd StepDefinition) Commands() []depVos.CommandDefinition {
	commandsCopy := make([]depVos.CommandDefinition, len(sd.commands))
	copy(commandsCopy, sd.commands)
	return commandsCopy
}

func (sd StepDefinition) Variables() []depVos.Variable {
	variablesCopy := make([]depVos.Variable, len(sd.variables))
	copy(variablesCopy, sd.variables)
	return variablesCopy
}

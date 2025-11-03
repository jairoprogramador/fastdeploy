package entities

import (
	"errors"
	"fmt"
	"slices"

	shared "github.com/jairoprogramador/fastdeploy-core/internal/domain/shared"

	defVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/vos"
)

type StepDefinition struct {
	name      string
	triggers  []defVos.TriggerDefinition
	commands  []defVos.CommandDefinition
	variables []defVos.VariableDefinition
}

func NewStepDefinition(
	name string,
	triggers []defVos.TriggerDefinition,
	commands []defVos.CommandDefinition,
	variables []defVos.VariableDefinition) (StepDefinition, error) {

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
		triggersCommon := []defVos.TriggerDefinition{
			defVos.TriggerFromString(shared.ScopeRecipe),
			defVos.TriggerFromString(shared.ScopeCode),
			defVos.TriggerFromString(shared.ScopeVars),
		}

		if name == shared.StepSupply {
			triggers = []defVos.TriggerDefinition{
				defVos.TriggerFromString(shared.ScopeRecipe),
				defVos.TriggerFromString(shared.ScopeVars),
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
			return StepDefinition{}, fmt.Errorf("variable %s duplicada", variable.Name())
		}
		variableNames[variable.Name()] = struct{}{}
	}

	triggersCopy := make([]defVos.TriggerDefinition, len(triggers))
	copy(triggersCopy, triggers)

	commandsCopy := make([]defVos.CommandDefinition, len(commands))
	copy(commandsCopy, commands)

	variablesCopy := make([]defVos.VariableDefinition, len(variables))
	copy(variablesCopy, variables)

	return StepDefinition{
		name:      name,
		triggers:  triggersCopy,
		commands:  commandsCopy,
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

func (sd StepDefinition) Commands() []defVos.CommandDefinition {
	commandsCopy := make([]defVos.CommandDefinition, len(sd.commands))
	copy(commandsCopy, sd.commands)
	return commandsCopy
}

func (sd StepDefinition) Variables() []defVos.VariableDefinition {
	variablesCopy := make([]defVos.VariableDefinition, len(sd.variables))
	copy(variablesCopy, sd.variables)
	return variablesCopy
}

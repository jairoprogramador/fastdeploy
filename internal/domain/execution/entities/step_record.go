package entities

import (
	"fmt"

	defEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/entities"

	exeSer "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/services"
	exeVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"
)

type StepRecord struct {
	name     string
	status   exeVos.StepStatus
	commands []*CommandRecord
}

func NewStepRecord(stepDef defEnt.StepDefinition) *StepRecord {
	commands := make([]*CommandRecord, 0, len(stepDef.Commands()))

	for _, cmdDef := range stepDef.Commands() {
		command := NewCommandRecord(cmdDef)
		commands = append(commands, command)
	}

	return &StepRecord{
		name:     stepDef.Name(),
		status:   exeVos.StepStatusPending,
		commands: commands,
	}
}

func (se *StepRecord) FinalizeCommand(
	commandName, command, commandOutput string,
	exitCode int,
	resolver exeSer.ResolverService,
) error {

	commandRecord, err := se.SearchCommand(commandName)
	if err != nil {
		return err
	}

	err = commandRecord.Finalize(command, commandOutput, exitCode, resolver)
	if err != nil {
		return err
	}

	se.updateStatus()

	return nil
}

func (se *StepRecord) SearchCommand(commandName string) (*CommandRecord, error) {
	for _, cmd := range se.commands {
		if cmd.Name() == commandName {
			return cmd, nil
		}
	}
	return nil, fmt.Errorf("no se encontró el comando '%s' en el paso '%s'", commandName, se.name)
}

func (se *StepRecord) updateStatus() {
	if se.status == exeVos.StepStatusSkipped {
		return
	}

	hasFailed := false
	allCompleted := true

	for _, cmd := range se.commands {
		if cmd.Status() == exeVos.CommandStatusFailed {
			hasFailed = true
			break
		}
		if cmd.Status() == exeVos.CommandStatusPending {
			allCompleted = false
			break
		}
	}

	if hasFailed {
		se.status = exeVos.StepStatusFailed
	} else if allCompleted {
		se.status = exeVos.StepStatusSuccessful
	} else {
		se.status = exeVos.StepStatusInProgress
	}
}

func (se *StepRecord) Name() string {
	return se.name
}

func (se *StepRecord) Status() exeVos.StepStatus {
	return se.status
}

func (se *StepRecord) Commands() []*CommandRecord {
	cmdsCopy := make([]*CommandRecord, len(se.commands))
	copy(cmdsCopy, se.commands)
	return cmdsCopy
}

func (se *StepRecord) Skip() {
	if se.status == exeVos.StepStatusPending {
		se.status = exeVos.StepStatusSkipped
	}
}

func (se *StepRecord) Cached() {
	if se.status == exeVos.StepStatusPending {
		se.status = exeVos.StepStatusCached
	}
}

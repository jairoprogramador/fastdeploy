package entities

import (
	"fmt"

	depEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/entities"

	orchSer "github.com/jairoprogramador/fastdeploy-core/internal/domain/orchestration/services"
	orchVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/orchestration/vos"
)

type StepRecord struct {
	name     string
	status   orchVos.StepStatus
	commands []*CommandRecord
}

func NewStepRecord(stepDef depEnt.StepDefinition) *StepRecord {
	commands := make([]*CommandRecord, 0, len(stepDef.Commands()))

	for _, cmdDef := range stepDef.Commands() {
		command := NewCommandRecord(cmdDef)
		commands = append(commands, command)
	}

	return &StepRecord{
		name:     stepDef.Name(),
		status:   orchVos.StepStatusPending,
		commands: commands,
	}
}

func (se *StepRecord) FinalizeCommand(
	commandName, commandResolved, record string,
	exitCode int,
	resolver orchSer.TemplateResolver,
) error {

	command, err := se.SearchCommand(commandName)
	if err != nil {
		return err
	}

	err = command.Finalize(commandResolved, record, exitCode, resolver)
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
	if se.status == orchVos.StepStatusSkipped {
		return
	}

	hasFailed := false
	allCompleted := true

	for _, cmd := range se.commands {
		if cmd.Status() == orchVos.CommandStatusFailed {
			hasFailed = true
			break
		}
		if cmd.Status() == orchVos.CommandStatusPending {
			allCompleted = false
			break
		}
	}

	if hasFailed {
		se.status = orchVos.StepStatusFailed
	} else if allCompleted {
		se.status = orchVos.StepStatusSuccessful
	} else {
		se.status = orchVos.StepStatusInProgress
	}
}

func (se *StepRecord) Name() string {
	return se.name
}

func (se *StepRecord) Status() orchVos.StepStatus {
	return se.status
}

func (se *StepRecord) Commands() []*CommandRecord {
	cmdsCopy := make([]*CommandRecord, len(se.commands))
	copy(cmdsCopy, se.commands)
	return cmdsCopy
}

func (se *StepRecord) Skip() {
	if se.status == orchVos.StepStatusPending {
		se.status = orchVos.StepStatusSkipped
	}
}

func (se *StepRecord) Cached() {
	if se.status == orchVos.StepStatusPending {
		se.status = orchVos.StepStatusCached
	}
}

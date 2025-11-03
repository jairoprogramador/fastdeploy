package entities

import (
	"errors"

	defVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/vos"

	exeSer "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/services"
	exeVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"
)

type CommandRecord struct {
	name          string
	command       string
	workdir       string
	templateFiles []string
	outputs       []exeVos.Output
	status        exeVos.CommandStatus
}

func NewCommandRecord(def defVos.CommandDefinition) *CommandRecord {

	outputs := make([]exeVos.Output, 0, len(def.Outputs()))
	for _, outputDef := range def.Outputs() {
		outputs = append(outputs, exeVos.NewOutput(outputDef))
	}

	return &CommandRecord{
		name:          def.Name(),
		command:       def.Cmd(),
		workdir:       def.Workdir(),
		templateFiles: def.TemplateFiles(),
		outputs:       outputs,
		status:        exeVos.CommandStatusPending,
	}
}

func (cr *CommandRecord) Finalize(
	command, commandOutput string,
	exitCode int,
	resolver exeSer.ResolverService) (err error) {

	if cr.status != exeVos.CommandStatusPending {
		return errors.New("solo se puede ejecutar un comando que está en estado pendiente")
	}

	cr.command = command

	if exitCode != 0 {
		cr.status = exeVos.CommandStatusFailed
		return nil
	}

	outputsExtracted, err := cr.extractOutputs(commandOutput, resolver)
	if err != nil {
		return err
	}

	cr.outputs = outputsExtracted
	cr.status = exeVos.CommandStatusSuccessful
	return nil
}

func (cr *CommandRecord) extractOutputs(
	record string, resolver exeSer.ResolverService) ([]exeVos.Output, error) {

	if len(cr.Outputs()) == 0 {
		return []exeVos.Output{}, nil
	}

	var outputs []exeVos.Output
	for _, output := range cr.Outputs() {
		outputExtracted, match, err := resolver.ResolveOutput(output, record)
		if err != nil || !match {
			cr.status = exeVos.CommandStatusFailed
			return nil, err
		}

		if output.Name() != "" {
			outputs = append(outputs, outputExtracted)
		}
	}
	return outputs, nil
}

func (cr *CommandRecord) Name() string {
	return cr.name
}

func (cr *CommandRecord) Command() string {
	return cr.command
}

func (cr *CommandRecord) Workdir() string {
	return cr.workdir
}

func (cr *CommandRecord) TemplateFiles() []string {
	return cr.templateFiles
}

func (cr *CommandRecord) Status() exeVos.CommandStatus {
	return cr.status
}

func (cr *CommandRecord) Outputs() []exeVos.Output {
	outputsCopy := make([]exeVos.Output, len(cr.outputs))
	copy(outputsCopy, cr.outputs)
	return outputsCopy
}

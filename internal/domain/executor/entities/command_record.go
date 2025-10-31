package entities

import (
	"errors"

	depVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/vos"

	orchSer "github.com/jairoprogramador/fastdeploy-core/internal/domain/executor/services"
	orchVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/executor/vos"
)

type CommandRecord struct {
	name          string
	command       string
	workdir       string
	templateFiles []string
	outputs       []orchVos.Output
	record        string
	status        orchVos.CommandStatus
}

func NewCommandRecord(def depVos.CommandDefinition) *CommandRecord {

	outputs := make([]orchVos.Output, 0, len(def.Outputs()))
	for _, outputDef := range def.Outputs() {
		outputs = append(outputs, orchVos.NewOutput(outputDef))
	}

	return &CommandRecord{
		name:          def.Name(),
		command:       def.Cmd(),
		workdir:       def.Workdir(),
		templateFiles: def.TemplateFiles(),
		outputs:       outputs,
		status:        orchVos.CommandStatusPending,
	}
}

func (cr *CommandRecord) Finalize(
	command, record string,
	exitCode int,
	resolver orchSer.TemplateResolver) (err error) {

	if cr.status != orchVos.CommandStatusPending {
		return errors.New("solo se puede ejecutar un comando que está en estado pendiente")
	}

	cr.command = command
	cr.record = record

	if exitCode != 0 {
		cr.status = orchVos.CommandStatusFailed
		return nil
	}

	outputsExtracted, err := cr.extractOutputs(record, resolver)
	if err != nil {
		return err
	}

	cr.outputs = outputsExtracted
	cr.status = orchVos.CommandStatusSuccessful
	return nil
}

func (cr *CommandRecord) extractOutputs(
	record string, resolver orchSer.TemplateResolver) ([]orchVos.Output, error) {

	if len(cr.Outputs()) == 0 {
		return []orchVos.Output{}, nil
	}

	var outputs []orchVos.Output
	for _, output := range cr.Outputs() {
		outputExtracted, match, err := resolver.ResolveOutput(output, record)
		if err != nil || !match {
			cr.status = orchVos.CommandStatusFailed
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

func (cr *CommandRecord) Status() orchVos.CommandStatus {
	return cr.status
}

func (cr *CommandRecord) Record() string {
	return cr.record
}

func (cr *CommandRecord) Outputs() []orchVos.Output {
	outputsCopy := make([]orchVos.Output, len(cr.outputs))
	copy(outputsCopy, cr.outputs)
	return outputsCopy
}

package aggregates

import (
	"fmt"

	defAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/aggregates"
	defEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/entities"

	exeEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/entities"
	exeSer "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/services"
	exeVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"
)

const (
	OutputOrderIdKey         = "order_id"
	OutputStepWorkdirKey     = "step_workdir"
	OutputCommWorkdirKey     = "commans_workdir"
	OutputProjectRevisionKey = "project_revision"
)

type ExecutionRecord struct {
	id            exeVos.ExecutionID
	status        exeVos.ExecutionStatus
	environment   string
	steps         []*exeEnt.StepRecord
	outputsShared map[string]exeVos.Output
}

func NewExecutionRecord(
	orderId exeVos.ExecutionID,
	deployment *defAgg.Deployment,
	environment string,
	finalStepName string,
	skippedStepNames map[string]struct{},
	initialOutputs []exeVos.Output,
) (*ExecutionRecord, error) {

	stepsRecords := getStepRecordFromDeployment(deployment.Steps(), skippedStepNames, finalStepName)

	outputsShared := getOutputsInitial(initialOutputs)

	newOrder := &ExecutionRecord{
		id:            orderId,
		status:        exeVos.ExecutionStatusInProgress,
		environment:   environment,
		steps:         stepsRecords,
		outputsShared: outputsShared,
	}

	newOrder.AddOutput(OutputOrderIdKey, orderId.String())
	return newOrder, nil
}

func getOutputsInitial(initialOutputs []exeVos.Output) map[string]exeVos.Output {
	outputsShared := make(map[string]exeVos.Output)
	for _, output := range initialOutputs {
		outputsShared[output.Name()] = output
	}
	return outputsShared
}

func getStepRecordFromDeployment(
	stepsDef []defEnt.StepDefinition,
	skippedStepNames map[string]struct{},
	finalStepName string) []*exeEnt.StepRecord {

	var stepsRecords []*exeEnt.StepRecord
	for _, stepDef := range stepsDef {
		stepRecord := exeEnt.NewStepRecord(stepDef)

		if _, shouldSkip := skippedStepNames[stepDef.Name()]; shouldSkip {
			stepRecord.Skip()
		}
		stepsRecords = append(stepsRecords, stepRecord)

		if stepDef.Name() == finalStepName {
			break
		}
	}
	return stepsRecords
}

func (o *ExecutionRecord) SearchStep(stepName string) (*exeEnt.StepRecord, error) {
	for _, step := range o.steps {
		if step.Name() == stepName {
			return step, nil
		}
	}
	return nil, fmt.Errorf("no se encontró el paso '%s' en la orden", stepName)
}

func (o *ExecutionRecord) FinalizeCommand(
	stepName, commandName, command, commandOutput string,
	exitCode int,
	resolver exeSer.ResolverService,
) error {
	stepRecord, err := o.SearchStep(stepName)
	if err != nil {
		return err
	}

	err = stepRecord.FinalizeCommand(commandName, command, commandOutput, exitCode, resolver)
	if err != nil {
		return err
	}

	commandRecord, err := stepRecord.SearchCommand(commandName)
	if err != nil {
		return err
	}

	for _, variable := range commandRecord.Outputs() {
		o.outputsShared[variable.Name()] = variable
	}

	o.updateStatus()

	return nil
}

func (o *ExecutionRecord) MarkStepAsCached(stepName string) {
	if stepRecord, err := o.SearchStep(stepName); err == nil {
		stepRecord.Cached()
		o.updateStatus()
	}
}

func (o *ExecutionRecord) updateStatus() {
	hasFailed := false
	allCompleted := true

	for _, step := range o.steps {
		if step.Status() == exeVos.StepStatusFailed {
			hasFailed = true
			break
		}
		if step.Status() == exeVos.StepStatusPending ||
			step.Status() == exeVos.StepStatusInProgress {
			allCompleted = false
			break
		}
	}

	if hasFailed {
		o.status = exeVos.ExecutionStatusFailed
	} else if allCompleted {
		o.status = exeVos.ExecutionStatusSuccessful
	} else {
		o.status = exeVos.ExecutionStatusInProgress
	}
}

func (o *ExecutionRecord) ID() exeVos.ExecutionID {
	return o.id
}

func (o *ExecutionRecord) Status() exeVos.ExecutionStatus {
	return o.status
}

func (o *ExecutionRecord) Environment() string {
	return o.environment
}

func (o *ExecutionRecord) StepsRecord() []*exeEnt.StepRecord {
	return o.steps
}

func (o *ExecutionRecord) Outputs() map[string]exeVos.Output {
	return o.outputsShared
}

func (o *ExecutionRecord) AddOutput(key, value string) {
	if variable, err := exeVos.NewOutputFromNameAndValue(key, value); err == nil {
		o.outputsShared[variable.Name()] = variable
	}
}

func (o *ExecutionRecord) AddOutputsMap(mapVariables map[string]string) {
	for key, value := range mapVariables {
		o.AddOutput(key, value)
	}
}

func (o *ExecutionRecord) RemoveOutput(key string) {
	delete(o.outputsShared, key)
}

func (o *ExecutionRecord) GetOutputsMapForSave() map[string]string {
	varsMap := make(map[string]string)
	for _, value := range o.outputsShared {
		varsMap[value.Name()] = value.Value()
	}
	return varsMap
}

func (o *ExecutionRecord) GetOutputsMapForFingerprint() map[string]string {
	varsMap := make(map[string]string)
	for _, value := range o.outputsShared {
		if value.Name() != OutputOrderIdKey &&
			value.Name() != OutputStepWorkdirKey &&
			value.Name() != OutputCommWorkdirKey &&
			value.Name() != OutputProjectRevisionKey {
			varsMap[value.Name()] = value.Value()
		}
	}
	return varsMap
}

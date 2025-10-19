package aggregates

import (
	"fmt"

	depAgg "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/aggregates"
	depEnt "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entities"
	depVos "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"

	orchEnt "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/entities"
	orchSer "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/services"
	orchVos "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"
)

const (
	OutputOrderIdName = "order_id"
)

type Order struct {
	id            orchVos.OrderID
	status        orchVos.OrderStatus
	environment   string
	steps         []*orchEnt.StepRecord
	outputsShared map[string]orchVos.Output
}

func NewOrder(
	orderId orchVos.OrderID,
	template *depAgg.DeploymentTemplate,
	environment depVos.Environment,
	finalStepName string,
	skippedStepNames map[string]struct{},
	initialOutputs []orchVos.Output,
) (*Order, error) {

	stepsRecords := getStepRecordFromTemplate(template.Steps(), skippedStepNames, finalStepName)

	outputsShared := getOutputsInitial(initialOutputs, orderId)

	return &Order{
		id:            orderId,
		status:        orchVos.OrderStatusInProgress,
		environment:   environment.Value(),
		steps:         stepsRecords,
		outputsShared: outputsShared,
	}, nil
}

func getOutputsInitial(initialOutputs []orchVos.Output, orderId orchVos.OrderID) map[string]orchVos.Output {
	outputsShared := make(map[string]orchVos.Output)
	for _, output := range initialOutputs {
		outputsShared[output.Name()] = output
	}
	outputOrderId, _ := orchVos.NewOutputFromNameAndValue(OutputOrderIdName, orderId.String())
	outputsShared[outputOrderId.Name()] = outputOrderId

	return outputsShared
}

func getStepRecordFromTemplate(
	stepsDef []depEnt.StepDefinition,
	skippedStepNames map[string]struct{},
	finalStepName string) []*orchEnt.StepRecord {

	var stepsRecords []*orchEnt.StepRecord
	for _, stepDef := range stepsDef {
		stepRecord := orchEnt.NewStepRecord(stepDef)

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

func (o *Order) SearchStep(stepName string) (*orchEnt.StepRecord, error) {
	for _, step := range o.steps {
		if step.Name() == stepName {
			return step, nil
		}
	}
	return nil, fmt.Errorf("no se encontró el paso '%s' en la orden", stepName)
}

func (o *Order) FinalizeCommand(
	stepName, commandName, commandResolved, record string,
	exitCode int,
	resolver orchSer.TemplateResolver,
) error {
	stepRecord, err := o.SearchStep(stepName)
	if err != nil {
		return err
	}

	err = stepRecord.FinalizeCommand(commandName, commandResolved, record, exitCode, resolver)
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

func (o *Order) MarkStepAsCached(stepName string) {
	if stepRecord, err := o.SearchStep(stepName); err == nil {
		stepRecord.Cached()
		o.updateStatus()
	}
}

func (o *Order) updateStatus() {
	hasFailed := false
	allCompleted := true

	for _, step := range o.steps {
		if step.Status() == orchVos.StepStatusFailed {
			hasFailed = true
			break
		}
		if step.Status() == orchVos.StepStatusPending ||
			step.Status() == orchVos.StepStatusInProgress {
			allCompleted = false
			break
		}
	}

	if hasFailed {
		o.status = orchVos.OrderStatusFailed
	} else if allCompleted {
		o.status = orchVos.OrderStatusSuccessful
	} else {
		o.status = orchVos.OrderStatusInProgress
	}
}

func (o *Order) ID() orchVos.OrderID {
	return o.id
}

func (o *Order) Status() orchVos.OrderStatus {
	return o.status
}

func (o *Order) Environment() string {
	return o.environment
}

func (o *Order) StepsRecord() []*orchEnt.StepRecord {
	return o.steps
}

func (o *Order) Outputs() map[string]orchVos.Output {
	return o.outputsShared
}

func (o *Order) AddOutput(key, value string) {
	variable, _ := orchVos.NewOutputFromNameAndValue(key, value)
	o.outputsShared[variable.Name()] = variable
}

func (o *Order) AddOutputsMap(mapVariables map[string]string) {
	for key, value := range mapVariables {
		o.AddOutput(key, value)
	}
}

func (o *Order) RemoveOutput(key string) {
	delete(o.outputsShared, key)
}

func (o *Order) GetOutputMap() map[string]string {
	varsMap := make(map[string]string)
	for _, value := range o.outputsShared {
		varsMap[value.Name()] = value.Value()
	}
	return varsMap
}

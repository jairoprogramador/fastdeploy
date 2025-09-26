package deployment

import (
	"slices"

	"github.com/jairoprogramador/fastdeploy/internal/domain/step/values"
)

var stepsOrder = []string{
	values.StepTest,
	values.StepSupply,
	values.StepPackage,
	values.StepDeploy,
}

type Orchestrator interface {
	CreateStepChain() (StepChain, error)
}

type OrchestratorImpl struct {
	targetStep string
	blockedSteps []string
}

func NewOrchestrator(targetStep string, blockedSteps []string) Orchestrator {
	return &OrchestratorImpl{
		targetStep: targetStep,
		blockedSteps: blockedSteps,
	}
}

func (f *OrchestratorImpl) CreateStepChain() (StepChain, error) {

	var firstCommand StepChain
	var lastCommand StepChain

	addCommand := func(c StepChain) {
		if firstCommand == nil {
			firstCommand = c
			lastCommand = c
		} else {
			lastCommand.SetNext(c)
			lastCommand = c
		}
	}

	for _, stepName := range stepsOrder {
		if !slices.Contains(f.blockedSteps, stepName) {
			stepChain := NewStepChain(stepName)
			addCommand(stepChain)
		}

		if stepName == f.targetStep {
			break
		}
	}

	return firstCommand, nil
}

package mapper

import (
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/executionstate/aggregates"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/executionstate/vos"
)

func StateStepsToDTO(stateSteps aggregates.StateSteps) *map[string]bool {
	stepMap := make(map[string]bool)
	for _, value := range stateSteps.GetStateSteps() {
		stepMap[value.GetName()] = value.IsSuccessful()
	}
	return &stepMap
}

func StateStepsToDomain(stepMap map[string]bool) (aggregates.StateSteps, error) {
	stateSteps := aggregates.NewStateSteps()
	for key, value := range stepMap {
		step, err := vos.NewStateStep(key, value)
		if err != nil {
			return aggregates.NewStateSteps(), err
		}
		stateSteps.AddStep(step)
	}
	return stateSteps, nil
}
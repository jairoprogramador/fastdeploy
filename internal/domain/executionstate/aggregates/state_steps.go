package aggregates

import "github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/vos"

type StateSteps struct {
	steps map[string]vos.StateStep
}

func NewStateSteps() StateSteps {
	return StateSteps{steps: make(map[string]vos.StateStep)}
}

func (s StateSteps) AddStep(step vos.StateStep) {
	s.steps[step.GetName()] = step
}

func (s StateSteps) GetStateSteps() map[string]vos.StateStep {
	stepsCopy := make(map[string]vos.StateStep, len(s.steps))
	for name, step := range s.steps {
		stepsCopy[name] = step
	}
	return stepsCopy
}

func (s StateSteps) IsStepAlreadyExecuted(name string) bool {
	step, ok := s.steps[name]
	if !ok {
		return false
	}
	return step.IsSuccessful()
}
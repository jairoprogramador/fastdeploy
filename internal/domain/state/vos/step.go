package vos

import "errors"

const (
	StepTest    = "test"
	StepSupply  = "supply"
	StepPackage = "package"
	StepDeploy  = "deploy"
)

type Step struct {
	value string
}

func NewStep(value string) (Step, error) {
	if value == "" {
		return Step{}, errors.New("step value cannot be empty")
	}
	return Step{value: value}, nil
}

func (s Step) String() string {
	return s.value
}

func (s Step) Equals(other Step) bool {
	return s.value == other.value
}

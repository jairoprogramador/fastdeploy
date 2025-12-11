package services

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/ports"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/services/matchers"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"
)

func NewStateMatcherFactory(step vos.Step, policy vos.CachePolicy) (ports.StateMatcher, error) {
	switch step.String() {
	case vos.StepTest:
		return &matchers.TestStateMatcher{Policy: policy}, nil
	case vos.StepSupply:
		return &matchers.SupplyStateMatcher{}, nil
	case vos.StepPackage:
		return &matchers.PackageStateMatcher{}, nil
	case vos.StepDeploy:
		return &matchers.DeployStateMatcher{}, nil
	default:
		return nil, fmt.Errorf("no state matcher found for step: %s", step)
	}
}

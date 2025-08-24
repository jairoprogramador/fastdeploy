package strategy

import "github.com/jairoprogramador/fastdeploy/internal/domain/deployment"

type StepStrategy interface {
	Execute(deployment.Context) error
}

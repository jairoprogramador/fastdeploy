package strategy

import "github.com/jairoprogramador/fastdeploy/internal/domain/context/service"

type StepStrategy interface {
	Execute(service.Context) error
}

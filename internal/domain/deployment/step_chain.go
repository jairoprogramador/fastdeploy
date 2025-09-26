package deployment

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/context/values"
	"github.com/jairoprogramador/fastdeploy/internal/domain/step/service"
)

type StepChain interface {
	SetNext(StepChain)
	Execute(stepService service.StepService, ctx *values.ContextValue) error
}

type StepChainImpl struct {
	next StepChain
	stepName string
}

func NewStepChain(stepName string) StepChain {
	return &StepChainImpl{
		stepName: stepName,
	}
}

func (b *StepChainImpl) SetNext(nextStep StepChain) {
	b.next = nextStep
}

func (b *StepChainImpl) ExecuteNext(stepService service.StepService, ctx *values.ContextValue) error {
	if b.next != nil {
		return b.next.Execute(stepService, ctx)
	}
	return nil
}

func (t *StepChainImpl) Execute(stepService service.StepService, ctx *values.ContextValue) error {
	step, err := stepService.Load(t.stepName)
	if err != nil {
		return err
	}

	if err := stepService.Run(step, ctx); err != nil {
		return err
	}
	return t.ExecuteNext(stepService, ctx)
}

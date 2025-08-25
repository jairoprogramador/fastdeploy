package deployment

import (
	"github.com/jairoprogramador/fastdeploy/internal/application/project"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/service"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
)

type ExecuteStep struct {
	readerProject   project.Reader
	context          deployment.Context
	stepOrchestrator service.StepOrchestrator
}

func NewExecuteStep(
	readerProject project.Reader,
	context deployment.Context,
	stepOrchestrator service.StepOrchestrator) *ExecuteStep {
	return &ExecuteStep{
		readerProject:    readerProject,
		context:          context,
		stepOrchestrator: stepOrchestrator,
	}
}

func (e *ExecuteStep) StartStep(stepName string, blockedSteps []string) error {
	project, err := e.readerProject.Read()
	if err != nil {
		return err
	}

	commandChain, err := e.stepOrchestrator.GetExecutionPlan(stepName, blockedSteps)
	if err != nil {
		return err
	}

	e.context.Set(constants.KeyTechnologyName, project.GetTechnology().GetName().Value())
	e.context.Set(constants.KeyTechnologyVersion, project.GetTechnology().GetVersion().Value())
	e.context.Set(constants.KeyNameRepository, project.GetRepository().GetURL().ExtractNameRepository())

	return commandChain.Execute(e.context)
}

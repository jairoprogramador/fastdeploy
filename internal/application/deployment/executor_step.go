package deployment

import (
	"github.com/jairoprogramador/fastdeploy/internal/application/project"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/service"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
)

type ExecuteStep struct {
	readerProject    project.Reader
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

	orchestrator, err := e.stepOrchestrator.GetExecutionPlan(stepName, blockedSteps)
	if err != nil {
		return err
	}

	e.context.Set(constants.KeyIdProject, project.GetID().Value())
	e.context.Set(constants.KeyNameProject, project.GetName().Value())
	e.context.Set(constants.KeyNameOrganization, project.GetOrganization().Value())
	e.context.Set(constants.KeyNameTeam, project.GetTeam().Value())
	e.context.Set(constants.KeyUrlRepository, project.GetRepository().GetURL().Value())
	e.context.Set(constants.KeyVersionRepository, project.GetRepository().GetVersion().Value())
	e.context.Set(constants.KeyNameRepository, project.GetRepository().GetURL().ExtractNameRepository())
	e.context.Set(constants.KeyNameTechnology, project.GetTechnology().Value())
	e.context.Set(constants.KeyVersionDeployment, project.GetDeployment().GetVersion().Value())

	return orchestrator.Execute(e.context)
}

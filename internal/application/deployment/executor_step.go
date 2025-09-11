package deployment

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/application/project"
	contextRepository "github.com/jairoprogramador/fastdeploy/internal/domain/context/port"
	contextService "github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
	deploymentService "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/service"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
)

type ExecuteStep struct {
	readerProject    project.Reader
	context          contextService.Context
	contextRepository contextRepository.Repository
	stepOrchestrator deploymentService.StepOrchestrator
}

func NewExecuteStep(
	readerProject project.Reader,
	context contextService.Context,
	contextRepository contextRepository.Repository,
	stepOrchestrator deploymentService.StepOrchestrator) *ExecuteStep {
	return &ExecuteStep{
		readerProject:    readerProject,
		context:          context,
		contextRepository: contextRepository,
		stepOrchestrator: stepOrchestrator,
	}
}

func (e *ExecuteStep) StartStep(stepName string, blockedSteps []string) error {
	project, err := e.readerProject.Read()
	if err != nil {
		return err
	}

	dataContext, err := e.contextRepository.Load(project.GetName().Value())
	if err != nil {
		return err
	}

	e.context.SetAll(dataContext.GetAll())

	orchestrator, err := e.stepOrchestrator.GetExecutionPlan(stepName, blockedSteps)
	if err != nil {
		return err
	}

	pathProject, err := e.readerProject.PathDirectory()
	if err != nil {
		return err
	}
	
	pathDeployment, err := e.readerProject.PathDirectoryGit(project)
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
	e.context.Set(constants.KeyPathProject, pathProject)
	e.context.Set(constants.KeyPathDeployment, pathDeployment)
	e.context.Set(constants.KeyEnvironmentName, "dev")
	e.context.Set(constants.KeySubscriptionId, "ee6f0101-cf12-48ca-b7b8-1745af77d759")

	fmt.Println("INICIO Contexto de la ejecución")
	for id, value := range e.context.GetAll() {
		fmt.Println(id,":", value)
	}

	err = orchestrator.Execute(e.context)
	if err != nil {
		return err
	}

	result := e.contextRepository.Save(project.GetName().Value(), e.context)
	if result != nil {
		return result
	}

	fmt.Println("FINAL Contexto de la ejecución")
	for id, value := range e.context.GetAll() {
		fmt.Println(id,":", value)
	}

	return result
}

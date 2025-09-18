package deployment

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/application/project"
	contextRepository "github.com/jairoprogramador/fastdeploy/internal/domain/context/port"
	contextService "github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
	deploymentService "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/port"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
)

type ExecuteStep struct {
	readerProject     project.Reader
	identifier        port.Identifier
	context           contextService.Context
	contextRepository contextRepository.Repository
	stepOrchestrator  deploymentService.StepOrchestrator
	validateEnvironment deploymentService.ValidateEnvironment
}

func NewExecuteStep(
	readerProject project.Reader,
	identifier port.Identifier,
	context contextService.Context,
	contextRepository contextRepository.Repository,
	stepOrchestrator deploymentService.StepOrchestrator,
	validateEnvironment deploymentService.ValidateEnvironment) *ExecuteStep {
	return &ExecuteStep{
		readerProject:     readerProject,
		identifier:        identifier,
		context:           context,
		contextRepository: contextRepository,
		stepOrchestrator:  stepOrchestrator,
		validateEnvironment: validateEnvironment,
	}
}

func (e *ExecuteStep) StartStep(stepName string, blockedSteps []string) error {

	project, err := e.readerProject.Read()
	if err != nil {
		return err
	}

	environment, err := e.context.Get(constants.Environment)
	if err != nil {
		return err
	}

	repositoryName := project.GetRepository().GetURL().ExtractNameRepository()

	valid, err := e.validateEnvironment.IsValidEnvironment(repositoryName, environment)
	if err != nil {
		return err
	}
	if !valid {
		return fmt.Errorf("environment %s is not valid", environment)
	}


	dataContext, err := e.contextRepository.Load(project.GetName().Value(), environment)
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

	deploymentId := e.identifier.Generate(project.GetName().Value(), project.GetDeployment().GetVersion().Value())

	e.context.Set(constants.ProjectId, project.GetID().Value()[0:4])
	e.context.Set(constants.ProjectId8, project.GetID().Value()[0:8])
	e.context.Set(constants.ProjectId12, project.GetID().Value()[0:12])
	e.context.Set(constants.ProjectId16, project.GetID().Value()[0:16])
	e.context.Set(constants.ProjectName, project.GetName().Value()[0:min(len(project.GetName().Value()), 4)])
	e.context.Set(constants.ProjectName8, project.GetName().Value()[0:min(len(project.GetName().Value()), 8)])
	e.context.Set(constants.ProjectName12, project.GetName().Value()[0:min(len(project.GetName().Value()), 12)])
	e.context.Set(constants.ProjectName16, project.GetName().Value()[0:min(len(project.GetName().Value()), 16)])
	e.context.Set(constants.ProjectOrganization, project.GetOrganization().Value())
	e.context.Set(constants.ProjectTeam, project.GetTeam().Value())
	e.context.Set(constants.ProjectCategory, project.GetCategory().Value())
	e.context.Set(constants.DeploymentRepositoryUrl, project.GetRepository().GetURL().Value())
	e.context.Set(constants.DeploymentRepositoryVersion, project.GetRepository().GetVersion().Value())
	e.context.Set(constants.DeploymentRepositoryName, repositoryName)
	e.context.Set(constants.ProjectTechnology, project.GetTechnology().Value())
	e.context.Set(constants.ProjectVersion, project.GetDeployment().GetVersion().Value())
	e.context.Set(constants.ProjectSourcePath, pathProject)
	e.context.Set(constants.DeploymentRepositoryPath, pathDeployment)
	e.context.Set(constants.Environment, environment)
	e.context.Set(constants.Environment8, environment[0:min(len(environment), 8)])
	e.context.Set(constants.DeploymentId, deploymentId[0:4])
	e.context.Set(constants.DeploymentId8, deploymentId[0:8])
	e.context.Set(constants.DeploymentId12, deploymentId[0:12])
	e.context.Set(constants.DeploymentId16, deploymentId[0:16])
	e.context.Set(constants.ToolName, constants.ToolName)
	e.context.Set(constants.Step, stepName)

	//fmt.Println("INICIO Contexto de la ejecución")
	//for id, value := range e.context.GetAll() {
	//fmt.Println(id,":", value)
	//}

	err = orchestrator.Execute(e.context)
	if err != nil {
		return err
	}

	result := e.contextRepository.Save(project.GetName().Value(), e.context)
	if result != nil {
		return result
	}

	//fmt.Println("FINAL Contexto de la ejecución")
	//for id, value := range e.context.GetAll() {
	//fmt.Println(id,":", value)
	//}

	return result
}

package deployment

import (
	"github.com/jairoprogramador/fastdeploy/internal/application/project"
	contextValue "github.com/jairoprogramador/fastdeploy/internal/domain/context/values"
	serviceRouter "github.com/jairoprogramador/fastdeploy/internal/domain/router/service"
	serviceDeployment "github.com/jairoprogramador/fastdeploy/internal/domain/deployment"
	serviceStep "github.com/jairoprogramador/fastdeploy/internal/domain/step/service"
	serviceShared "github.com/jairoprogramador/fastdeploy/internal/domain/shared/service"
)

type CommandExecutor struct {
	readerProject project.Reader
	routerHome serviceRouter.HomeRouterService
	stepService serviceStep.StepService
}

func NewCommandExecutor(
	readerProject project.Reader,
	routerHome serviceRouter.HomeRouterService,
	routerProject serviceRouter.ProjectRouterService,
	stepService serviceStep.StepService,
	) *CommandExecutor {
	return &CommandExecutor{
		readerProject: readerProject,
		routerHome: routerHome,
		stepService: stepService,
	}
}

func (e *CommandExecutor) ExecuteCommand(environment string, targetStep string, blockedSteps []string) error {

	pathFastDeploy, err := e.routerHome.GetPathFastDeploy()
	if err != nil {
		return err
	}

	pathWorkdirProject, err := e.routerHome.GetPathWorkdir()
	if err != nil {
		return err
	}

	project, err := e.readerProject.Read()
	if err != nil {
		return err
	}

	deploymentId := serviceShared.Generate(
		project.GetName().Value(), project.GetDeployment().GetVersion().Value())

	context, err := contextValue.NewContext(
		&project, environment, pathFastDeploy, deploymentId, pathWorkdirProject)
	if err != nil {
		return err
	}

	orchestrator := serviceDeployment.NewOrchestrator(targetStep, blockedSteps)
	stepChain, err := orchestrator.CreateStepChain()
	if err != nil {
		return err
	}

	//fmt.Println("INICIO Contexto de la ejecución")
	//for id, value := range e.context.GetAll() {
	//fmt.Println(id,":", value)
	//}


	//fmt.Println("FINAL Contexto de la ejecución")
	//for id, value := range e.context.GetAll() {
	//fmt.Println(id,":", value)
	//}

	return stepChain.Execute(e.stepService, context)
}

package main

import (
	//"log"

	app "github.com/jairoprogramador/fastdeploy/internal/application/deployment"
	/* projectApp "github.com/jairoprogramador/fastdeploy/internal/application/project"
	portCommand "github.com/jairoprogramador/fastdeploy/internal/domain/command/port"
	serviceCommand "github.com/jairoprogramador/fastdeploy/internal/domain/command/service"
	portContext "github.com/jairoprogramador/fastdeploy/internal/domain/context/port"
	serviceContext "github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
	portRouter "github.com/jairoprogramador/fastdeploy/internal/domain/router/port"
	routerService "github.com/jairoprogramador/fastdeploy/internal/domain/router/service"
	values "github.com/jairoprogramador/fastdeploy/internal/domain/router/values"
	portStep "github.com/jairoprogramador/fastdeploy/internal/domain/step/port"
	stepService "github.com/jairoprogramador/fastdeploy/internal/domain/step/service"
	portVariable "github.com/jairoprogramador/fastdeploy/internal/domain/variable/port"
	serviceVariable "github.com/jairoprogramador/fastdeploy/internal/domain/variable/service"
	deploymentService "github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/service"
	projectService "github.com/jairoprogramador/fastdeploy/internal/infrastructure/project/service" */
)

func GetEnvironmentRepository() []string {
	/* repositoryProject := projectService.NewFileRepository()
	readerProject := projectApp.NewReader(repositoryProject)
	proj, err := readerProject.Read()
	if err != nil {
		log.Printf("Advertencia: no se ha podido leer el proyecto para crear subcomandos de deploy: %v", err)
	}
	repoName := proj.GetRepository().GetURL().ExtractNameRepository()
	environments, err := deploymentService.NewEnvironmentRepository().GetEnvironments(repoName)
	if err != nil {
		log.Printf("Advertencia: no se pudieron obtener los entornos: %v", err)
	}
	return environments */
	return []string{"local", "development", "production"}
}

func GetCommandExecutor() *app.CommandExecutor {
	/* project, err := readerProject.Read()
	if err != nil {
		log.Printf("Advertencia: no se ha podido leer el proyecto: %v", err)
	} */

	/* portHomeRouter := portRouter.HomeRouter{}
	homeRouterService := routerService.NewHomeRouterService(portHomeRouter)

	portProjectRouter := portRouter.ProjectRouter{}
	newParameter, err := values.NewParameter("pathFastDeploy", "projectName", "repositoryName", "environment", "stackName")
	if err != nil {
		log.Printf("Advertencia: no se ha podido crear el parametro para el router de proyecto: %v", err)
	}
	projectRouterService := routerService.NewProjectRouterService(portProjectRouter, newParameter)
	contextPort := portContext.ContextPort{}

	contextService := serviceContext.NewContextService(projectRouterService, contextPort)

	portRepositoryRouter := portRouter.RepositoryRouter{}
	repoRouterService := routerService.NewRepositoryRouterService(portRepositoryRouter, newParameter)

	portVariable := portVariable.VariablePort{}
	variableService := serviceVariable.NewVariableService(portVariable, repoRouterService) */

	//contextValue, err := contextValue.NewContext(&project, "environment", "homeDir", "deploymentId", "workdirProject")
	/* if err != nil {
		log.Printf("Advertencia: no se ha podido crear el contexto: %v", err)
	} */
	/* commandService := serviceCommand.NewCommandService(
		portCommand.TemplatePort{},
		portCommand.ExecutorPort{},
		portCommand.WorkdirPort{},
		variableService,
		repoRouterService,
		projectRouterService)

	stepPort := portStep.StepPort{}
	stepService := stepService.NewStepService(
		contextService, variableService, commandService, repoRouterService, stepPort)

	repositoryProject := projectService.NewFileRepository()
	readerProject := projectApp.NewReader(repositoryProject) */

	//return app.NewCommandExecutor(readerProject, homeRouterService, projectRouterService, stepService)
	return nil
}

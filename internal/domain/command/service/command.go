package service

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/command/port"
	command "github.com/jairoprogramador/fastdeploy/internal/domain/command/values"
	context "github.com/jairoprogramador/fastdeploy/internal/domain/context/values"
	serviceRouter "github.com/jairoprogramador/fastdeploy/internal/domain/router/service"
	serviceVariable "github.com/jairoprogramador/fastdeploy/internal/domain/variable/service"
	variable "github.com/jairoprogramador/fastdeploy/internal/domain/variable/values"
)

type CommandService interface {
	Execute(
		command command.CommandValue,
		stepName string,
		context *context.ContextValue) ([]variable.VariableValue, error)
}

type CommandServiceImpl struct {
	templatePort     port.TemplatePort
	executorPort     port.ExecutorPort
	workdirPort      port.WorkdirPort
	variableService  serviceVariable.VariableService
	routerRepository serviceRouter.RepositoryRouterService
	routerProject    serviceRouter.ProjectRouterService
	processorLine    func(line string, context *context.ContextValue) string
}

func NewCommandService(
	templateCommand port.TemplatePort,
	executorCommand port.ExecutorPort,
	workdirCommand port.WorkdirPort,
	variableService serviceVariable.VariableService,
	repositoryRouter serviceRouter.RepositoryRouterService,
	projectRouter serviceRouter.ProjectRouterService) CommandService {

	processorLine := func(line string, context *context.ContextValue) string {
		return variableService.Process(line, context)
	}

	return &CommandServiceImpl{
		templatePort:     templateCommand,
		executorPort:     executorCommand,
		workdirPort:      workdirCommand,
		variableService:  variableService,
		routerRepository: repositoryRouter,
		routerProject:    projectRouter,
		processorLine:    processorLine,
	}
}

func (c *CommandServiceImpl) Execute(
	commandValue command.CommandValue,
	stepName string,
	context *context.ContextValue) ([]variable.VariableValue, error) {

	workdir := commandValue.GetWorkdir()

	if workdir != command.DEFAULT_WORKDIR {
		stepPathRepository := c.routerRepository.GetPathStep(stepName)
		stepPathProject := c.routerProject.GetPathStep(stepName)

		if err := c.workdirPort.Copy(stepPathRepository, stepPathProject); err != nil {
			return []variable.VariableValue{}, err
		}
		workdir = c.routerRepository.BuildPath(stepPathProject, commandValue.GetWorkdir())
	}

	if err := c.processTemplatePath(commandValue, stepName, context); err != nil {
		return []variable.VariableValue{}, err
	}

	command := c.processorLine(commandValue.GetCommand(), context)

	outputCommand, err := c.executorPort.Run(command, workdir)

	variables, err := commandValue.IsValid(outputCommand, err)
	if err != nil {
		return []variable.VariableValue{}, err
	}

	return variables, nil
}

func (c *CommandServiceImpl) processTemplatePath(
	command command.CommandValue,
	stepName string,
	context *context.ContextValue) error {

	if len(command.GetTemplates()) == 0 {
		return nil
	}

	for _, templatePath := range command.GetTemplates() {

		stepPathRepository := c.routerRepository.GetPathStep(stepName)
		fullPathRepository := c.routerRepository.BuildPath(stepPathRepository, command.GetWorkdir(), templatePath)

		stepPathProject := c.routerProject.GetPathStep(stepName)
		fullPathProject := c.routerRepository.BuildPath(stepPathProject, command.GetWorkdir(), templatePath)

		if err := c.workdirPort.Copy(fullPathRepository, fullPathProject); err != nil {
			return err
		}

		if err := c.templatePort.Process(fullPathProject, c.processorLine, context); err != nil {
			return err
		}
	}

	return nil
}

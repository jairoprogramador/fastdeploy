package service

import (
	valueStep "github.com/jairoprogramador/fastdeploy/internal/domain/step/values"
	valueContext "github.com/jairoprogramador/fastdeploy/internal/domain/context/values"
	serviceContext "github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
	serviceVariable "github.com/jairoprogramador/fastdeploy/internal/domain/variable/service"
	serviceCommand "github.com/jairoprogramador/fastdeploy/internal/domain/command/service"
	serviceRouter "github.com/jairoprogramador/fastdeploy/internal/domain/router/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/step/port"
)

const STEP_COMMANDS_FILE_NAME = "commands.yaml"

type StepService interface {
	Run(step valueStep.StepValue, context *valueContext.ContextValue) error
	Load(stepName string) (valueStep.StepValue, error)
}

type StepServiceImpl struct {
	contextService serviceContext.ContextService
	variableService serviceVariable.VariableService
	commandService serviceCommand.CommandService
	routerRepositoryService serviceRouter.RepositoryRouterService
	stepPort port.StepPort
}

func NewStepService(
	contextService serviceContext.ContextService,
	variableService serviceVariable.VariableService,
	commandService serviceCommand.CommandService,
	routerService serviceRouter.RepositoryRouterService,
	stepPort port.StepPort) StepService {

	return &StepServiceImpl{
		contextService: contextService,
		variableService: variableService,
		commandService: commandService,
		routerRepositoryService: routerService,
		stepPort: stepPort,
	}
}

func (s *StepServiceImpl) Run(
	step valueStep.StepValue,
	context *valueContext.ContextValue) error {

	contextStep, err := s.contextService.AddContextStep(step.GetName(), context)

	if err != nil {
		return err
	}

	contextStep, err = s.variableService.AddVariablesComputed(step.GetName(), contextStep)
	if err != nil {
		return err
	}

	contextStep, err = s.variableService.AddVariablesStep(step.GetName(), contextStep)
	if err != nil {
		return err
	}

	for _, command := range step.GetCommands() {
		variables, err := s.commandService.Execute(command, step.GetName(), contextStep)
		if err != nil {
			return err
		}

		for _, variable := range variables {
			contextStep.Set(variable.GetName(), variable.GetValue())
		}
	}

	err = s.contextService.SaveContextStep(step.GetName(), contextStep)
	if err != nil {
		return err
	}

	return err
}

func (s *StepServiceImpl) Load(stepName string) (valueStep.StepValue, error) {
	pathStep := s.routerRepositoryService.GetPathStep(stepName)
	pathStepFile := s.routerRepositoryService.BuildPath(pathStep, STEP_COMMANDS_FILE_NAME)
	commands, err := s.stepPort.LoadCommands(pathStepFile)
	if err != nil {
		return valueStep.StepValue{}, err
	}
	return valueStep.NewStepValue(stepName, commands)
}
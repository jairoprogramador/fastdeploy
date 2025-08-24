package deployment

import (
	"github.com/jairoprogramador/fastdeploy/internal/application/project"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/service"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
)

type ExecuteStep struct {
	readerProject  project.Reader
	context        deployment.Context
	commandManager service.CommandManager
}

func NewExecuteStep(
	readerProject project.Reader,
	context deployment.Context,
	commandManager service.CommandManager) *ExecuteStep {
	return &ExecuteStep{
		readerProject:  readerProject,
		context:        context,
		commandManager: commandManager,
	}
}

func (e *ExecuteStep) StartStep(stepName string, blockedSteps []string) error {
	project, err := e.readerProject.Read()
	if err != nil {
		return err
	}

	strategy, err := e.commandManager.CreateChain(stepName, blockedSteps)
	if err != nil {
		return err
	}

	e.context.Set(constants.KeyTechnologyName, project.GetTechnology().GetName().Value())
	e.context.Set(constants.KeyTechnologyVersion, project.GetTechnology().GetVersion().Value())
	e.context.Set(constants.KeyNameRepository, project.GetRepository().GetURL().ExtractNameRepository())

	return strategy.Execute(e.context)
}

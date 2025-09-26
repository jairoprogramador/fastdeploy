package service

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/context/port"
	"github.com/jairoprogramador/fastdeploy/internal/domain/context/values"
	"github.com/jairoprogramador/fastdeploy/internal/domain/router/service"
)

const CONTEXT_FILE_NAME = "context.gob"

type ContextService interface {
	AddContextStep(stepName string, contextValue *values.ContextValue) (*values.ContextValue, error)
	SaveContextStep(stepName string, contextValue *values.ContextValue) error
}

type ContextServiceImpl struct {
	projectRouter service.ProjectRouterService
	contextPort   port.ContextPort
}

func NewContextService(
	projectRouter service.ProjectRouterService,
	contextPort port.ContextPort) ContextService {
	return &ContextServiceImpl{
		projectRouter: projectRouter,
		contextPort:   contextPort,
	}
}

func (c *ContextServiceImpl) AddContextStep(
	stepName string,
	contextValue *values.ContextValue) (*values.ContextValue, error) {

	pathFileContextStep := c.getPathFileContextStep(stepName)

	data, err := c.contextPort.Load(pathFileContextStep)
	if err != nil {
		return contextValue, err
	}

	contextValue.AddAll(data)

	return contextValue, nil
}

func (c *ContextServiceImpl) SaveContextStep(stepName string, context *values.ContextValue) error {
	pathFileContextStep := c.getPathFileContextStep(stepName)
	return c.contextPort.Save(pathFileContextStep, context.GetAll())
}

func (c *ContextServiceImpl) getPathFileContextStep(stepName string) string {
	pathStep := c.projectRouter.GetPathStep(stepName)
	pathFile := c.projectRouter.BuildPath(pathStep, CONTEXT_FILE_NAME)
	return pathFile
}
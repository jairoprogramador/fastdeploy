package values

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entity"
)

const (
	ToolName            = "fastdeploy"
	DeploymentId        = "deployment_id"
	DeploymentId8       = "deployment_id8"
	DeploymentId12      = "deployment_id12"
	DeploymentId16      = "deployment_id16"
	ProjectOrganization = "project_organization"
	ProjectTeam         = "project_team"
	ProjectCategory     = "project_category"
	ProjectName         = "project_name"
	ProjectName8        = "project_name8"
	ProjectName12       = "project_name12"
	ProjectName16       = "project_name16"
	ProjectId           = "project_id"
	ProjectId8          = "project_id8"
	ProjectId12         = "project_id12"
	ProjectId16         = "project_id16"
	ProjectVersion      = "project_version"
	ProjectTechnology   = "project_technology"
	ProjectSourcePath   = "project_source_path"
	ProjectType         = "project_type"
	Environment         = "environment"
	Environment4        = "environment4"
	Environment8        = "environment8"
)

/*
type Context interface {
	Get(key string) (string, error)
	Set(key, value string)
	GetAll() map[string]string
	SetAll(data map[string]string)
	GetProject() entity.Project
	GetProjectName() string
	//GetRepositoryName() string
	GetEnvironment() string
	GetHomeDir() string
} */

const DEFAULT_ENVIRONMENT = "local"

type ContextValue struct {
	mu          sync.RWMutex
	variables   map[string]string
	project     entity.Project
	environment string
	homeDirFastDeploy     string
	deploymentId string
	workdirProject string
}

func NewContext(
	project *entity.Project,
	environment, homeDirFastdeploy,
	deploymentId, workdirProject string) (*ContextValue, error) {

	environment = strings.TrimSpace(environment)
	homeDirFastdeploy = strings.TrimSpace(homeDirFastdeploy)
	workdirProject = strings.TrimSpace(workdirProject)

	if project == nil {
		return &ContextValue{}, errors.New("project cannot be empty")
	}

	if environment == "" {
		environment = DEFAULT_ENVIRONMENT
	}

	if homeDirFastdeploy == "" {
		return &ContextValue{}, errors.New("homeDirFastdeploy cannot be empty")
	}

	if workdirProject == "" {
		return &ContextValue{}, errors.New("workdirProject cannot be empty")
	}

	context := &ContextValue{
		variables:   make(map[string]string),
		project:     *project,
		environment: environment,
		homeDirFastDeploy:     homeDirFastdeploy,
		deploymentId: deploymentId,
		workdirProject: workdirProject,
	}

	setVariablesDefault(context)

	return context, nil
}

func (c *ContextValue) Get(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	value, ok := c.variables[key]
	if !ok {
		return "", fmt.Errorf("variable no encontrada: %s", key)
	}
	return value, nil
}

func (c *ContextValue) Set(key, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.variables[key] = value
}

func (c *ContextValue) GetAll() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.variables
}

func (c *ContextValue) SetAll(data map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.variables = data
}

func (c *ContextValue) AddAll(data map[string]string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, value := range data {
		c.variables[key] = value
	}
}

func (c *ContextValue) GetProject() entity.Project {
	return c.project
}

func (c *ContextValue) GetProjectName() string {
	return c.project.GetName().Value()
}

func (c *ContextValue) GetProjectId() string {
	return c.project.GetID().Value()
}

func (c *ContextValue) GetProjectVersion() string {
	return c.project.GetDeployment().GetVersion().Value()
}

func (c *ContextValue) GetProjectTechnology() string {
	return c.project.GetTechnology().Value()
}

func (c *ContextValue) GetEnvironment() string {
	return c.environment
}

func (c *ContextValue) GetHomeDirFastDeploy() string {
	return c.homeDirFastDeploy
}

func (c *ContextValue) GetDeploymentId() string {
	return c.deploymentId
}

func (c *ContextValue) GetWorkdirProject() string {
	return c.workdirProject
}

func setVariablesDefault(context *ContextValue) {
	context.Set(ToolName, ToolName)
	context.Set(DeploymentId, getSubstringDeploymentId(context, 4))
	context.Set(DeploymentId8, getSubstringDeploymentId(context, 8))
	context.Set(DeploymentId12, getSubstringDeploymentId(context, 12))
	context.Set(DeploymentId16, getSubstringDeploymentId(context, 16))
	context.Set(ProjectName, getSubstringProjectName(context, 4))
	context.Set(ProjectName8, getSubstringProjectName(context, 8))
	context.Set(ProjectName12, getSubstringProjectName(context, 12))
	context.Set(ProjectName16, getSubstringProjectName(context, 16))
	context.Set(ProjectId, getSubstringProjectId(context, 4))
	context.Set(ProjectId8, getSubstringProjectId(context, 8))
	context.Set(ProjectId12, getSubstringProjectId(context, 12))
	context.Set(ProjectId16, getSubstringProjectId(context, 16))
	context.Set(ProjectVersion, context.GetProjectVersion())
	context.Set(ProjectTechnology, context.GetProjectTechnology())
	context.Set(ProjectSourcePath, context.GetWorkdirProject())
	context.Set(ProjectType, context.GetProject().GetCategory().Value())
	context.Set(ProjectCategory, context.GetProject().GetCategory().Value())
	context.Set(ProjectOrganization, context.GetProject().GetOrganization().Value())
	context.Set(ProjectTeam, context.GetProject().GetTeam().Value())
	context.Set(Environment, context.GetEnvironment())
	context.Set(Environment4, getSubstring(context.GetEnvironment(), 4))
	context.Set(Environment8, getSubstring(context.GetEnvironment(), 8))
}

func getSubstringDeploymentId(context *ContextValue, length int) string {
	return getSubstring(context.GetDeploymentId(), length)
}

func getSubstringProjectName(context *ContextValue, length int) string {
	return getSubstring(context.GetProjectName(), length)
}

func getSubstringProjectId(context *ContextValue, length int) string {
	return getSubstring(context.GetProjectId(), length)
}

func getSubstring(value string, length int) string {
	return value[0:min(len(value), length)]
}
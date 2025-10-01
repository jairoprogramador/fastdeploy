package dto

import (
	"context"

	deploymentvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	orchestrationvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
)

// ExecuteOrderRequest es el DTO (Data Transfer Object) que encapsula todos los
// parámetros necesarios para ejecutar una orden.
type ExecuteOrderRequest struct {
	Ctx              context.Context
	TemplateSource   deploymentvos.TemplateSource
	EnvironmentName  string
	FinalStepName    string
	ProjectName      string
	ProjectRootPath  string
	SkippedStepNames map[string]struct{}
	InitialVariables []orchestrationvos.Variable
}
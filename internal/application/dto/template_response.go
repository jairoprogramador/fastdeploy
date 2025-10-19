package dto

import "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/aggregates"

type TemplateResponse struct {
	Template       *aggregates.DeploymentTemplate
	TemplatePath   string
	RepositoryName string
}

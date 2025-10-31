package dto

import "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/aggregates"

type TemplateResponse struct {
	Template       *aggregates.DeploymentTemplate
	TemplatePath   string
	RepositoryName string
}

package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/dom/dto"
)

func DomToDTO(dom *aggregates.DeploymentObjectModel) dto.DomDTO {
	return dto.DomDTO{
		Product:    ProductToDTO(dom.Product()),
		Project:    ProjectToDTO(dom.Project()),
		Template:   TemplateToDTO(dom.Template()),
		Technology: TechnologyToDTO(dom.Technology()),
	}
}

func DomToDomain(dto dto.DomDTO) (*aggregates.DeploymentObjectModel, error) {
	product, err := ProductToDomain(dto.Product)
	if err != nil {
		return nil, err
	}
	project, err := ProjectToDomain(dto.Project)
	if err != nil {
		return nil, err
	}
	template, err := TemplateToDomain(dto.Template)
	if err != nil {
		return nil, err
	}
	technology, err := TechnologyToDomain(dto.Technology)
	if err != nil {
		return nil, err
	}
	return aggregates.NewDeploymentObjectModel(
		product,
		project,
		template,
		technology,
	), nil
}

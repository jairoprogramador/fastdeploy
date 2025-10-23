// Package mapper provides functions to map between DTOs and domain models.
package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/dom/dto"
)

func ToDomain(fileConfig dto.FileConfig) (*aggregates.Config, error) {
	project, err := ProjectToDomain(fileConfig.Project)
	if err != nil {
		return nil, err
	}

	template := TemplateToDomain(fileConfig.Template)
	technology := TechnologyToDomain(fileConfig.Technology)
	runtime := RuntimeToDomain(fileConfig.Runtime)
	state := StateToDomain(fileConfig.State)

	return aggregates.NewConfig(
		project,
		template,
		technology,
		runtime,
		state,
	), nil
}

func ToDTO(config *aggregates.Config) dto.FileConfig {

	projectDTO := ProjectToDTO(config.Project())
	templateDTO := TemplateToDTO(config.Template())
	technologyDTO := TechnologyToDTO(config.Technology())
	runtimeDTO := RuntimeToDTO(config.Runtime())
	stateDTO := StateToDTO(config.State())

	return dto.FileConfig {
		Project: projectDTO,
		Template: templateDTO,
		Technology: technologyDTO,
		Runtime: runtimeDTO,
		State: stateDTO,
	}
}
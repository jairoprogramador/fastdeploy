package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/project/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/project/dto"
)

func ToDomain(fileConfig dto.FileConfig) (*aggregates.Config, error) {
	project, err := ProjectToDomain(fileConfig.Project)
	if err != nil {
		return nil, err
	}

	template := TemplateToDomain(fileConfig.Template)
	state := StateToDomain(fileConfig.State)

	return aggregates.NewConfig(
		project,
		template,
		state,
	), nil
}

func ToDTO(config *aggregates.Config) dto.FileConfig {

	projectDTO := ProjectToDTO(config.Project())
	templateDTO := TemplateToDTO(config.Template())
	stateDTO := StateToDTO(config.State())

	return dto.FileConfig {
		Project: projectDTO,
		Template: templateDTO,
		State: stateDTO,
	}
}
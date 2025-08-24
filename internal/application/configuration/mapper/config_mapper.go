package mapper

import (
	"github.com/jairoprogramador/fastdeploy/internal/application/configuration/dto"
	"github.com/jairoprogramador/fastdeploy/internal/domain/configuration/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/values"
)

func ToDomain(dto dto.ConfigDto) (entity.Configuration, error) {
	organization, err := values.NewNameOrganization(dto.NameOrganization)
	if err != nil {
		return entity.Configuration{}, err
	}

	team, err := values.NewTeam(dto.Team)
	if err != nil {
		return entity.Configuration{}, err
	}

	repositoryUrl, err := values.NewUrlRepository(dto.UrlRepository)
	if err != nil {
		return entity.Configuration{}, err
	}

	repository := values.NewRepository(repositoryUrl)

	technologyName, err := values.NewNameTechnology(dto.Technology.Name)
	if err != nil {
		return entity.Configuration{}, err
	}

	technologyVersion, err := values.NewVersionTechnology(dto.Technology.Version)
	if err != nil {
		return entity.Configuration{}, err
	}

	technology := values.NewTechnology(technologyName, technologyVersion)

	return entity.NewConfiguration(organization, team, repository, technology), nil
}
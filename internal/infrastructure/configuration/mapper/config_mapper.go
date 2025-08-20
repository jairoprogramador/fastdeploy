package mapper

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/configuration/entities"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/values"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/configuration/dto"
)

func ToDto(configuration entities.Configuration) (dto.ConfigDto, error) {
	organization := configuration.GetNameOrganization().Value()
	team := configuration.GetTeam().Value()
	repository := configuration.GetRepository().GetURL().Value()

	technology := dto.TechnologyInfo{
		Name:    configuration.GetTechnology().GetName().Value(),
		Version: configuration.GetTechnology().GetVersion().Value(),
	}

	return dto.ConfigDto{
		Organization: organization,
		Team:         team,
		Repository:   repository,
		Technology:   technology,
	}, nil
}

func ToDomain(dto dto.ConfigDto) (entities.Configuration, error) {
	organization, err := values.NewNameOrganization(dto.Organization)
	if err != nil {
		return entities.Configuration{}, err
	}

	team, err := values.NewTeam(dto.Team)
	if err != nil {
		return entities.Configuration{}, err
	}

	repositoryUrl, err := values.NewUrlRepository(dto.Repository)
	if err != nil {
		return entities.Configuration{}, err
	}
	repository := values.NewRepository(repositoryUrl)

	technologyName, err := values.NewNameTechnology(dto.Technology.Name)
	if err != nil {
		return entities.Configuration{}, err
	}

	technologyVersion, err := values.NewVersionTechnology(dto.Technology.Version)
	if err != nil {
		return entities.Configuration{}, err
	}

	technology := values.NewTechnology(technologyName, technologyVersion)

	return entities.NewConfiguration(
		organization,
		team,
		repository,
		technology,
	), nil
}

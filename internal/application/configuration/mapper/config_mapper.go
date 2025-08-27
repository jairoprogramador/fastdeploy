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

	repository := values.NewRepository(repositoryUrl, values.NewDefaultVersionRepository())

	technology, err := values.NewNameTechnology(dto.Technology)
	if err != nil {
		return entity.Configuration{}, err
	}

	return entity.NewConfiguration(organization, team, repository, technology), nil
}
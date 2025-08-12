package utils

import (
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/config"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/project"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/repository"
)

func GetRepositoryPath(projectEntity project.ProjectEntity) (string, error) {
	configDirPath, err := config.GetConfigDirPath()
	if err != nil {
		return "", err
	}

	return repository.GetRepositoryDirPath(configDirPath, projectEntity.Repository), nil
}

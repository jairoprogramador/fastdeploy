package command

import (
	"deploy/internal/infrastructure/repository"
)

func IsInitialized() bool {
	projectRepository := repository.GetProjectRepository()
	return projectRepository.Exists()
}
package command

import (
	"deploy/internal/infrastructure/repository"
	"deploy/internal/domain/service"
)

func IsInitialized() bool {
	projectRepository := repository.GetProjectRepository()
	projectService := service.GetProjectService(projectRepository)

	return projectService.IsInitialized()
}
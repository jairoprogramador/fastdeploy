package command

import (
	"deploy/internal/infrastructure/repository"
	"deploy/internal/domain/service"
)

func IsInitialized() bool {
	globalConfigRepository := repository.NewGlobalConfigRepository()
	projectRepository := repository.NewProjectRepository()

	globalConfigService := service.NewGlobalConfigService(globalConfigRepository)
	projectService := service.NewProjectService(projectRepository, globalConfigService)
	
	_, err := projectService.Load()
	return err == nil
}
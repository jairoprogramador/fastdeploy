package command

import (
	"deploy/internal/infrastructure/repository"
	"deploy/internal/domain/service"
	"deploy/internal/application/dto"
)

func IsInitialized() *dto.ResponseDto {
	globalConfigRepository := repository.GetGlobalConfigRepository()
	projectRepository := repository.GetProjectRepository()

	globalConfigService := service.GetGlobalConfigService(globalConfigRepository)
	fileRepository := repository.GetFileRepository()
	projectService := service.GetProjectService(projectRepository, fileRepository, globalConfigService)
	
	_, err := projectService.Load()
	return dto.GetDtoWithError(err)
}
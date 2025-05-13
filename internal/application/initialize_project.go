package command

import (
	"deploy/internal/infrastructure/repository"
	"deploy/internal/application/dto"
	"deploy/internal/domain/service"
)

func InitializeProject() *dto.ResponseDto {
	globalConfigRepository := repository.GetGlobalConfigRepository()
	projectRepository := repository.GetProjectRepository()

	globalConfigService := service.GetGlobalConfigService(globalConfigRepository)
	projectService := service.GetProjectService(projectRepository, globalConfigService)

	return dto.GetDtoWithModel(projectService.Initialize())
}

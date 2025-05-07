package command

import (
	"deploy/internal/infrastructure/repository"
	"deploy/internal/application/dto"
	"deploy/internal/domain/service"
)

func InitializeProject() *dto.ResponseDto {
	globalConfigRepository := repository.NewGlobalConfigRepository()
	projectRepository := repository.NewProjectRepository()

	globalConfigService := service.NewGlobalConfigService(globalConfigRepository)
	projectService := service.NewProjectService(projectRepository, globalConfigService)

	return dto.GetNewResponseDtoFromModel(projectService.Initialize())
}

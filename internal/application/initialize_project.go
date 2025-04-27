package command

import (
	"deploy/internal/infrastructure/repository"
	"deploy/internal/application/dto"
	"deploy/internal/domain/service"
)

func InitializeProject() *dto.ResponseDto {
	projectRepository := repository.GetProjectRepository()
	projectService := service.GetProjectService(projectRepository)

	initializeService := service.GetInitializeService(*projectService)

	resp := initializeService.Initialize()
	return dto.GetNewResponseDtoFromModel(resp)
}

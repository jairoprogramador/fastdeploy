package command

import (
	"deploy/internal/infrastructure/repository"
	"deploy/internal/application/dto"
	"deploy/internal/domain/service"
	"deploy/internal/domain"
)

func InitializeProject() *dto.Message  {
	projectRepository := repository.GetProjectRepository()
	projectService := service.GetProjectService(projectRepository)

	isInitialized := projectService.IsInitialized()

	if isInitialized {
		return dto.GetNewMessage(constants.MessagePreviouslyInitializedProject)
	}

	if resp := projectService.CreateProyect(); resp.Error != nil {
        return dto.GetNewMessageFromResponse(resp)
    }

	if resp:= projectService.CreateDockerfileTemplate(); resp.Error != nil {
        return dto.GetNewMessageFromResponse(resp)
    }

	if resp:= projectService.CreateDockercomposeTemplate(); resp.Error != nil {
        return dto.GetNewMessageFromResponse(resp)
    }

	return dto.GetNewMessage(constants.MessageSuccessInitializingProject)
}

package command

import (
	"deploy/internal/infrastructure/repository"
	"deploy/internal/application/dto"
	"deploy/internal/domain/service"
)

func PublishProject() *dto.Message  {
	publishRepository := repository.GetPublishRepository()
	publishService := service.GetPublishService(publishRepository)

	responseBuild := publishService.Build()
	if responseBuild.Error != nil {
		return dto.GetNewMessageFromResponse(responseBuild) 
	}

	responsePackage := publishService.Package(responseBuild)
	if responsePackage.Error != nil {
		return dto.GetNewMessageFromResponse(responsePackage)
	}

	resp := publishService.Deliver(responsePackage)
	return dto.GetNewMessageFromResponse(resp)
}

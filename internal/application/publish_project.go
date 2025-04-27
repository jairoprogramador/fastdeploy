package command

import (
	"deploy/internal/infrastructure/repository"
	"deploy/internal/application/dto"
	"deploy/internal/domain/service"
)

func PublishProject() *dto.ResponseDto  {
	publishRepository := repository.GetPublishRepository()
	publishService := service.GetPublishService(publishRepository)

	responseBuild := publishService.Build()
	if responseBuild.Error != nil {
		return dto.GetNewResponseDtoFromModel(responseBuild) 
	}

	responsePackage := publishService.Package(responseBuild)
	if responsePackage.Error != nil {
		return dto.GetNewResponseDtoFromModel(responsePackage)
	}

	responseDeliver := publishService.Deliver(responsePackage)
	return dto.GetNewResponseDtoFromModel(responseDeliver)
}

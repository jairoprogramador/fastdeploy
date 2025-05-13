package command

import (
	"deploy/internal/infrastructure/repository"
	"deploy/internal/application/dto"
	"deploy/internal/domain/service"
)

func AddSonarQube() *dto.ResponseDto {
	sonarqubeRepository := repository.GetSonarqubeRepository()
	sonarqubeService := service.GetSonarqubeService(sonarqubeRepository)

	resp := sonarqubeService.Add()
	return dto.GetDtoWithModel(resp)
}
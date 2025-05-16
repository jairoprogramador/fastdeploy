package application

import (
	"deploy/internal/application/dto"
	"deploy/internal/domain/model"
)

var (
	projectModel *model.Project
)

func Initialize() *dto.ResponseDto {
	return dto.GetDtoWithModel(getProjectService().Initialize())
}

func IsInitialize() *dto.ResponseDto {
	var err error
	projectModel, err = getProjectService().Load()
	return dto.GetDtoWithError(err)
}

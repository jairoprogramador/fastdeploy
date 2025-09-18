package mapper

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entity"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/dto"
)

func ToDomainList(dtoList []dto.EnvironmentDto) ([]entity.Environment, error) {
	var environmentList []entity.Environment
	for _, dto := range dtoList {
		environmentList = append(environmentList, entity.NewEnvironment(dto.Value))
	}
	return environmentList, nil
}

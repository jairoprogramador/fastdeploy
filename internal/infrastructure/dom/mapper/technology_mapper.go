package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/vos"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/dom/dto"
)

func TechnologyToDomain(dto dto.TechnologyDTO) vos.Technology {
	return vos.NewTechnology(dto.Stack, dto.Infrastructure)
}

func TechnologyToDTO(technology vos.Technology) dto.TechnologyDTO {
	return dto.TechnologyDTO{
		Stack: technology.Stack(),
		Infrastructure: technology.Infrastructure(),
	}
}
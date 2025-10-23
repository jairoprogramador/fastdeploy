package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/vos"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/dom/dto"
)

func StateToDomain(dto dto.StateDTO) vos.State {
	return vos.NewState(dto.Backend, dto.URL)
}

func StateToDTO(state vos.State) dto.StateDTO {
	return dto.StateDTO{
		Backend: state.Backend(),
		URL: state.URL(),
	}
}
package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/vos"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/dom/dto"
)

func RuntimeToDomain(dto dto.RuntimeDTO) vos.Runtime {
	return vos.NewRuntime(
		vos.NewImage(dto.Image.Source, dto.Image.Tag),
		vos.NewVolumes(dto.Volumes.ProjectMountPath, dto.Volumes.StateMountPath),
	)
}

func RuntimeToDTO(runtime vos.Runtime) dto.RuntimeDTO {
	return dto.RuntimeDTO{
		Image: dto.ImageDTO{
			Source: runtime.Image().Source(),
			Tag: runtime.Image().Tag(),
		},
		Volumes: dto.VolumesDTO{
			ProjectMountPath: runtime.Volumes().ProjectMountPath(),
			StateMountPath: runtime.Volumes().StateMountPath(),
		},
	}
}
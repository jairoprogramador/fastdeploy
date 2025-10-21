package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/shared/vos"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/dom/dto"
)


func TemplateToDomain(dto dto.TemplateDTO) (vos.TemplateSource, error) {
	return vos.NewTemplateSource(
		dto.RepositoryURL,
		dto.Ref,
	)
}

func TemplateToDTO(template vos.TemplateSource) dto.TemplateDTO {
	return dto.TemplateDTO{
		RepositoryURL: template.Url(),
		Ref:           template.Ref(),
	}
}

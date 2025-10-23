package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/dom/vos"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/dom/dto"
)

func TemplateToDomain(dto dto.TemplateDTO) vos.Template {
	return vos.NewTemplate(dto.URL, dto.Ref)
}

func TemplateToDTO(template vos.Template) dto.TemplateDTO {
	return dto.TemplateDTO{
		URL: template.URL(),
		Ref: template.Ref(),
	}
}
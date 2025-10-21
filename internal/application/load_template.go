package application

import (
	"context"

	appDto "github.com/jairoprogramador/fastdeploy-core/internal/application/dto"
	depPor "github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/ports"
	shaVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/shared/vos"
)

type LoadTemplateService struct {
	templateRepository depPor.TemplateRepository
}

func NewLoadTemplateService(templateRepository depPor.TemplateRepository) *LoadTemplateService {
	return &LoadTemplateService{
		templateRepository: templateRepository,
	}
}

func (s *LoadTemplateService) Load(ctx context.Context, source shaVos.TemplateSource) (appDto.TemplateResponse, error) {
	template, templatePath, err := s.templateRepository.Load(ctx, source)
	if err != nil {
		return appDto.TemplateResponse{}, err
	}

	return appDto.TemplateResponse{
		Template:       template,
		TemplatePath:   templatePath,
		RepositoryName: source.NameTemplate(),
	}, nil
}

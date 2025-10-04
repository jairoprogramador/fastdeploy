package application

import (
	"context"

	"github.com/jairoprogramador/fastdeploy/newinternal/application/dto"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/ports"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
)

type LoadTemplateService struct {
	templateRepository ports.TemplateRepository
}

func NewLoadTemplateService(templateRepository ports.TemplateRepository) *LoadTemplateService {
	return &LoadTemplateService{
		templateRepository: templateRepository,
	}
}

func (s *LoadTemplateService) Load(ctx context.Context, repositoryURL string, ref string) (dto.LoadTemplateResponse, error) {
	templateSource, err := vos.NewTemplateSource(repositoryURL, ref)
	if err != nil {
		return dto.LoadTemplateResponse{}, err
	}
	template, templatePath, err := s.templateRepository.GetTemplate(ctx, templateSource)
	if err != nil {
		return dto.LoadTemplateResponse{}, err
	}
	repositoryName, err := s.templateRepository.GetRepositoryName(repositoryURL)
	if err != nil {
		return dto.LoadTemplateResponse{}, err
	}
	return dto.LoadTemplateResponse{
		Template: template,
		TemplatePath: templatePath,
		RepositoryName: repositoryName,
	}, nil
}
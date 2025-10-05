package mapper

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/dom/aggregates"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/dom/dto"
	"github.com/jairoprogramador/fastdeploy/internal/domain/dom/vos"
)

func ToDTO(dom *aggregates.DeploymentObjectModel) dto.DOMDTO {
	return dto.DOMDTO{
		Product:    toProductDTO(*dom.Product()),
		Project:    toProjectDTO(*dom.Project()),
		Template:   toTemplateDTO(*dom.Template()),
		Technology: toTechnologyDTO(*dom.Technology()),
	}
}

func ToDomain(dto dto.DOMDTO) (*aggregates.DeploymentObjectModel, error) {
	product, err := toProductDomain(dto.Product)
	if err != nil {
		return nil, err
	}
	project, err := toProjectDomain(dto.Project)
	if err != nil {
		return nil, err
	}
	template, err := toTemplateDomain(dto.Template)
	if err != nil {
		return nil, err
	}
	technology, err := toTechnologyDomain(dto.Technology)
	if err != nil {
		return nil, err
	}
	return aggregates.NewDeploymentObjectModel(
		product,
		project,
		template,
		technology,
	), nil
}

func toProductDomain(dto dto.ProductDTO) (*vos.Product, error) {
	productId := vos.ProductID(dto.ID)
	return vos.NewProduct(
		productId,
		dto.Name,
		dto.Description,
		dto.Team,
		dto.Organization,
	)
}

func toProjectDomain(dto dto.ProjectDTO) (*vos.Project, error) {
	projectId := vos.ProjectID(dto.ID)
	return vos.NewProject(
		projectId,
		dto.Name,
		dto.Version,
		dto.Description,
		dto.Team,
	)
}

func toTemplateDomain(dto dto.TemplateDTO) (*vos.Template, error) {
	return vos.NewTemplate(
		dto.RepositoryURL,
		dto.Ref,
	)
}

func toTechnologyDomain(dto dto.TechnologyDTO) (*vos.Technology, error) {
	return vos.NewTechnology(
		dto.Type,
		dto.Solution,
		dto.Stack,
		dto.Infrastructure,
	)
}

func toProductDTO(product vos.Product) dto.ProductDTO {
	return dto.ProductDTO{
		ID:   string(product.ID()),
		Name: product.Name(),
		Description: product.Description(),
		Team: product.Team(),
		Organization: product.Organization(),
	}
}

func toProjectDTO(project vos.Project) dto.ProjectDTO {
	return dto.ProjectDTO{
		ID:   string(project.ID()),
		Name: project.Name(),
		Version: project.Version(),
		Description: project.Description(),
		Team: project.Team(),
	}
}

func toTemplateDTO(template vos.Template) dto.TemplateDTO {
	return dto.TemplateDTO{
		RepositoryURL: template.RepositoryURL(),
		Ref: template.Ref(),
	}
}

func toTechnologyDTO(technology vos.Technology) dto.TechnologyDTO {
	return dto.TechnologyDTO{
		Type: technology.TypeTechnology(),
		Solution: technology.Solution(),
		Stack: technology.Stack(),
		Infrastructure: technology.Infrastructure(),
	}
}
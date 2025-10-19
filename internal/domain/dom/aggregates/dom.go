package aggregates

import (
	"fmt"

	domSer "github.com/jairoprogramador/fastdeploy/internal/domain/dom/services"
	domVos "github.com/jairoprogramador/fastdeploy/internal/domain/dom/vos"
	sharedVos "github.com/jairoprogramador/fastdeploy/internal/domain/shared/vos"
)

type DeploymentObjectModel struct {
	product    domVos.Product
	project    domVos.Project
	template   sharedVos.TemplateSource
	technology domVos.Technology
}

func NewDeploymentObjectModel(
	product domVos.Product,
	project domVos.Project,
	template sharedVos.TemplateSource,
	tech domVos.Technology) *DeploymentObjectModel {
	return &DeploymentObjectModel{
		product:    product,
		project:    project,
		template:   template,
		technology: tech,
	}
}

func (dom *DeploymentObjectModel) SetProjectRevision(revision string) {
	dom.project = dom.project.WithRevision(revision)
}

func (dom *DeploymentObjectModel) IsModified(shaGenerator domSer.ShaGenerator) bool {
	expectedProductID := shaGenerator.GenerateProductID(dom.product.Name(), dom.product.Organization())
	if dom.product.ID() != expectedProductID {
		return true
	}

	expectedProjectID := shaGenerator.GenerateProjectID(dom.technology)
	return dom.project.ID() != expectedProjectID
}

func (dom *DeploymentObjectModel) UpdateIDs(shaGenerator domSer.ShaGenerator) error {

	generatedProductId := shaGenerator.GenerateProductID(dom.product.Name(), dom.product.Organization())

	if dom.product.ID() != generatedProductId {
		newProduct, err := domVos.NewProduct(
			generatedProductId,
			dom.product.Name(),
			dom.product.Description(),
			dom.product.Team(),
			dom.product.Organization())
		if err != nil {
			return fmt.Errorf("failed to update product ID: %w", err)
		}
		dom.product = newProduct
	}

	generatedProjectId := shaGenerator.GenerateProjectID(dom.technology)

	if dom.project.ID() != generatedProjectId {
		newProject, err := domVos.NewProject(
			generatedProjectId,
			dom.project.Name(),
			dom.project.Version(),
			dom.project.Description(),
			dom.project.Team())
		if err != nil {
			return fmt.Errorf("failed to update project ID: %w", err)
		}
		dom.project = newProject
	}
	return nil
}

func (dom *DeploymentObjectModel) Product() domVos.Product               { return dom.product }
func (dom *DeploymentObjectModel) Project() domVos.Project               { return dom.project }
func (dom *DeploymentObjectModel) Template() sharedVos.TemplateSource { return dom.template }
func (dom *DeploymentObjectModel) Technology() domVos.Technology         { return dom.technology }

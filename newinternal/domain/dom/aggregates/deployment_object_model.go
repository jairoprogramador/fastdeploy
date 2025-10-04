package aggregates

import (
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/dom/services"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/dom/vos"
)

// DeploymentObjectModel es el Agregado Raíz para el contexto DOM.
// Encapsula y protege la consistencia de todo el archivo dom.yaml.
type DeploymentObjectModel struct {
	product    *vos.Product
	project    *vos.Project
	template   *vos.Template
	technology *vos.Technology
}

// NewDeploymentObjectModel es el constructor para crear un nuevo DOM.
func NewDeploymentObjectModel(
	product *vos.Product,
	project *vos.Project,
	template *vos.Template,
	tech *vos.Technology) *DeploymentObjectModel {
	return &DeploymentObjectModel{
		product:    product,
		project:    project,
		template:   template,
		technology: tech,
	}
}

// VerifyAndUpdateIDs es el método principal que encapsula la lógica de negocio.
// Comprueba si los IDs actuales son consistentes con los datos y los actualiza si es necesario.
// Devuelve 'true' si se realizó un cambio.
func (dom *DeploymentObjectModel) VerifyAndUpdateIDs(idGen services.IDGenerator) (bool, error) {
	isModified := false

	expectedProductID := idGen.GenerateProductID(dom.product.Name(), dom.product.Organization())
	if dom.product.ID() != expectedProductID {
		// Actualizar el producto con el nuevo ID
		updatedProduct, err := vos.NewProduct(
			expectedProductID,
			dom.product.Name(),
			dom.product.Description(),
			dom.product.Team(),
			dom.product.Organization())

		if err != nil {
			return false, err
		}
		dom.product = updatedProduct
		isModified = true
	}

	expectedProjectID := idGen.GenerateProjectID(*dom.technology)
	if dom.project.ID() != expectedProjectID {
		// Actualizar el proyecto con el nuevo ID
		updatedProject, err := vos.NewProject(
			expectedProjectID,
			dom.project.Name(),
			dom.project.Version(),
			dom.project.Description(),
			dom.project.Team())

		if err != nil {
			return false, err
		}
		dom.project = updatedProject
		isModified = true
	}

	return isModified, nil
}

// Getters para los VOs internos...
func (dom *DeploymentObjectModel) Product() *vos.Product       { return dom.product }
func (dom *DeploymentObjectModel) Project() *vos.Project       { return dom.project }
func (dom *DeploymentObjectModel) Template() *vos.Template     { return dom.template }
func (dom *DeploymentObjectModel) Technology() *vos.Technology { return dom.technology }

package services

import "github.com/jairoprogramador/fastdeploy/newinternal/domain/dom/vos"

// IDGenerator define el contrato para un servicio que puede generar
// los IDs consistentes para el agregado DOM.
type IDGenerator interface {
	GenerateProductID(name, organization string) vos.ProductID
	GenerateProjectID(tech vos.Technology) vos.ProjectID
}

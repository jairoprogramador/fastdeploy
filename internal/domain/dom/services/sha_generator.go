package services

import "github.com/jairoprogramador/fastdeploy/internal/domain/dom/vos"

type ShaGenerator interface {
	GenerateProductID(name, organization string) vos.ProductID
	GenerateProjectID(tech vos.Technology) vos.ProjectID
}

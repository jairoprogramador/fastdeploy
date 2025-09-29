package ports

import (
	"context"

	deploymententities "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/entities"
	deploymentvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
)

// StepVariableRepository define el contrato para un adaptador que puede cargar
// variables de configuración específicas para un paso en un ambiente determinado.
// Esto permite al dominio obtener toda la configuración necesaria sin conocer
// la fuente (e.g., archivos YAML, una API de configuración, etc.).
type StepVariableRepository interface {
	// Load obtiene una lista de variables para una combinación específica de
	// definición de paso y ambiente. La implementación se encargará de
	// encontrar y parsear el archivo correspondiente (e.g., variables/stag/supply.yaml).
	Load(
		ctx context.Context,
		repositoryName string,
		environment deploymentvos.Environment,
		stepDefinition deploymententities.StepDefinition,
	) ([]vos.Variable, error)
}

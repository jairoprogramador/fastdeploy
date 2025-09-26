package ports

import (
	"context"

	"github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/aggregates"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
)

// TemplateRepository define el contrato para obtener un agregado DeploymentTemplate.
// Actúa como un puerto en la arquitectura hexagonal, permitiendo que la capa de
// aplicación solicite un agregado sin conocer los detalles de su obtención (e.g., git, filesystem).
type TemplateRepository interface {
	// GetTemplate obtiene una plantilla de despliegue basada en su origen (URL y referencia).
	// La implementación se encargará de la lógica de clonado, checkout y parsing de los archivos
	// para construir un agregado válido y consistente.
	GetTemplate(ctx context.Context, source vos.TemplateSource) (*aggregates.DeploymentTemplate, error)
}

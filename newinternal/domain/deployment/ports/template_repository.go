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
	// GetTemplate obtiene una plantilla de despliegue y devuelve también la ruta local
	// al repositorio clonado para que otros servicios puedan usarla.
	GetTemplate(ctx context.Context, source vos.TemplateSource) (template *aggregates.DeploymentTemplate, repoLocalPath string, err error)

	GetRepositoryName(repoURL string) (string, error)
}

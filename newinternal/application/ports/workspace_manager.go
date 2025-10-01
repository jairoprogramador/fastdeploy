package ports

import (
	"context"

	deploymententities "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/entities"
	deploymentvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
)

// WorkspaceManager define el contrato para un adaptador que gestiona
// los directorios de trabajo para cada paso de la ejecución.
type WorkspaceManager interface {
	// PrepareStepWorkspace prepara y devuelve la ruta a un directorio de trabajo
	// limpio para un paso y ambiente específicos. La implementación se encarga de
	// la lógica de copiado de plantillas.
	PrepareStepWorkspace(
		ctx context.Context,
		projectName string,
		targetEnv deploymentvos.Environment,
		stepDef deploymententities.StepDefinition,
		templateRepoPath string,
	) (workspacePath string, err error)
}

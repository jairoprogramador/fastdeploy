package ports

// WorkspaceManager define el contrato para un adaptador que gestiona
// los directorios de trabajo para cada paso de la ejecución.
type WorkspaceManager interface {
	// PrepareStepWorkspace prepara y devuelve la ruta a un directorio de trabajo
	// limpio para un paso y ambiente específicos. La implementación se encarga de
	// la lógica de copiado de plantillas.
	PrepareStepWorkspace(
		projectName string,
		environment string,
		stepName string,
		repositoryName string,
	) (workspacePath string, err error)
}

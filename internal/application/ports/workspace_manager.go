package ports

type WorkspaceManager interface {
	Prepare(stepName string) (workspacePath string, err error)
}

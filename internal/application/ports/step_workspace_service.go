package ports

import "github.com/jairoprogramador/fastdeploy-core/internal/application/dto"

type StepWorkspaceService interface {
	Prepare(namesRequest dto.NamesParams, runRequest dto.RunParams) (workspacePath string, err error)
}

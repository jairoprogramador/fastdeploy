package ports

import (
	"context"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"
)

type CommandExecutor interface {
	Execute(
		ctx context.Context,
		command vos.Command,
		currentVars vos.VariableSet,
		workspaceStep string) *vos.ExecutionResult
}

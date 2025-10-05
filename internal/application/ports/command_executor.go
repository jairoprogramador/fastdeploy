package ports

import "context"

type CommandExecutor interface {
	Execute(ctx context.Context, workdir, command string) (log string, exitCode int, err error)
	CreateWorkDir(workdirs ...string) string
}

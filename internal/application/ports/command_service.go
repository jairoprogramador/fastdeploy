package ports

import "context"

type CommandService interface {
	Run(ctx context.Context, workdir, command string) (string, int, error)
	CreateWorkDir(workdirs ...string) string
}

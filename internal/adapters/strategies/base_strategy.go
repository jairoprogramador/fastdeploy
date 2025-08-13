package strategies

import "github.com/jairoprogramador/fastdeploy/internal/adapters/executor"

type BaseStrategy struct {
	RepositoryPath string
	Executor       executor.ExecutorCmd
}

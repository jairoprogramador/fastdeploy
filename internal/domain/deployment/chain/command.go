package chain

import "github.com/jairoprogramador/fastdeploy/internal/domain/deployment"

type CommandChain interface {
	SetNext(CommandChain)
	Execute(deployment.Context) error
}

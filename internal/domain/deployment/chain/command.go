package chain

import "github.com/jairoprogramador/fastdeploy/internal/domain/context/service"

type CommandChain interface {
	SetNext(CommandChain)
	Execute(service.Context) error
}

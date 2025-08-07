package commands

import "github.com/jairoprogramador/fastdeploy/internal/core/domain/context"

type Command interface {
	SetNext(Command)
	Execute(context.Context) error
}

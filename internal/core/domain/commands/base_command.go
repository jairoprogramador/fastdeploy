package commands

import "github.com/jairoprogramador/fastdeploy/internal/core/domain/context"

type BaseCommand struct {
	next Command
}

func (b *BaseCommand) SetNext(c Command) {
	b.next = c
}

func (b *BaseCommand) ExecuteNext(ctx context.Context) error {
	if b.next != nil {
		return b.next.Execute(ctx)
	}
	return nil
}

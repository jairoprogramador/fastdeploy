package chain

import "github.com/jairoprogramador/fastdeploy/internal/domain/deployment"

type BaseCommandChain struct {
	next CommandChain
}

func (b *BaseCommandChain) SetNext(nextCommand CommandChain) {
	b.next = nextCommand
}

func (b *BaseCommandChain) ExecuteNext(ctx deployment.Context) error {
	if b.next != nil {
		return b.next.Execute(ctx)
	}
	return nil
}

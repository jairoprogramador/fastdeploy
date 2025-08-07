package commands

type BaseCommand struct {
	next Command
}

func (b *BaseCommand) SetNext(c Command) {
	b.next = c
}

func (b *BaseCommand) ExecuteNext() error {
	if b.next != nil {
		return b.next.Execute()
	}
	return nil
}
package commands

type Command interface {
	SetNext(Command)
	Execute() error
}
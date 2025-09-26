package port

type ExecutorPort interface {
	Run(command string, workdir string) (string, error)
}

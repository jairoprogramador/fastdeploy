package repository

type VariableRepository interface {
	GetCommitHash() (string, error)
	GetCommitAuthor(commitHash string) (string, error)
	GetCommitMessage(commitHash string) (string, error)
}
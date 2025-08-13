package project

type RepositoryReader interface {
	Load() (*ProjectEntity, error)
	Exists() bool
}

type RepositoryWriter interface {
	Save(projectEntity ProjectEntity) error
	Delete() error
}

type ProjectRepository interface {
	RepositoryReader
	RepositoryWriter
}

package project

type ServiceReader interface {
	Load() (*ProjectEntity, error)
	Exists() bool
}

type ServiceWriter interface {
	Save(projectEntity ProjectEntity) error
	//Delete() error
}

type ProjectService interface {
	ServiceReader
	ServiceWriter
}

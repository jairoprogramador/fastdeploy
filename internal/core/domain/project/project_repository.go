package project

type ProjectRepository interface {
	Save(projectEntity ProjectEntity) error
	Load() (*ProjectEntity, error)
	Exists() bool
	Delete() error
}

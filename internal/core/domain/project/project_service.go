package project

type ProjectService interface {
	Save(projectEntity ProjectEntity) error
	Load() (*ProjectEntity, error)
	Exists() bool
	//Delete() error
}

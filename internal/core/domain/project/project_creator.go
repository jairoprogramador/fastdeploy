package project

type ProjectCreator interface {
	Create() (*ProjectEntity, error)
}

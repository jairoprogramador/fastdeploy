package project

type ProjectValidator interface {
	Validate(project ProjectEntity) error
}

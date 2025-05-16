package repository

type ContainerRepository interface {
	GetFullPathResource() (string, error)
	GetContentTemplate(pathTemplate string, params any) (string, error)	
}

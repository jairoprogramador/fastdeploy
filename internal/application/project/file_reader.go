package project

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entities"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/ports"
)

type Reader interface {
	ExistsFile() (bool, error)
	Read() (entities.Project, error)
}

type FileReader struct {
	repository ports.Repository
}

func NewReader(repository ports.Repository) Reader {
	return &FileReader{
		repository: repository,
	}
}

func (cs *FileReader) Read() (entities.Project, error) {
	project, err := cs.repository.Load()
	if err != nil {
		return entities.Project{}, err
	}

	return project, nil
}

func (cs *FileReader) ExistsFile() (bool, error) {
	return cs.repository.Exists()
}
package project

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/port"
)

type Reader interface {
	ExistsFile() (bool, error)
	Read() (entity.Project, error)
}

type FileReader struct {
	repository port.Repository
}

func NewReader(repository port.Repository) Reader {
	return &FileReader{
		repository: repository,
	}
}

func (cs *FileReader) Read() (entity.Project, error) {
	project, err := cs.repository.Load()
	if err != nil {
		return entity.Project{}, err
	}

	return project, nil
}

func (cs *FileReader) ExistsFile() (bool, error) {
	return cs.repository.Exists()
}
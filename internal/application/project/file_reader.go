package project

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/port"
)

type Reader interface {
	ExistsFile() (bool, error)
	Read() (entity.Project, error)
	PathDirectory() (string, error)
	PathDirectoryGit(project entity.Project) (string, error)
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

func (cs *FileReader) PathDirectory() (string, error) {
	return cs.repository.PathDirectory()
}

func (cs *FileReader) PathDirectoryGit(project entity.Project) (string, error) {
	return cs.repository.PathDirectoryGit(project)
}
package project

import (
	domainEntity "github.com/jairoprogramador/fastdeploy/internal/domain/project/entities"
	domainPort "github.com/jairoprogramador/fastdeploy/internal/domain/project/ports"
	"github.com/jairoprogramador/fastdeploy/internal/application/project/ports"
)

type FileReader struct {
	repository domainPort.Repository
}

func NewReader(repository domainPort.Repository) ports.Reader {
	return &FileReader{
		repository: repository,
	}
}

func (cs *FileReader) Read() (domainEntity.Project, error) {
	project, err := cs.repository.Load()
	if err != nil {
		return domainEntity.Project{}, err
	}

	return project, nil
}

func (cs *FileReader) ExistsFile() (bool, error) {
	return cs.repository.Exists()
}
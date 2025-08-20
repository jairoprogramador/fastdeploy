package project

import (
	"github.com/jairoprogramador/fastdeploy/internal/application/project/ports"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entities"
	domainPort "github.com/jairoprogramador/fastdeploy/internal/domain/project/ports"
)

type FileWriter struct {
	repository domainPort.Repository
}

func NewWriter(repository domainPort.Repository) ports.Writer {
	return &FileWriter{
		repository: repository,
	}
}

func (cs *FileWriter) Write(project entities.Project) error {
	if err := cs.repository.Save(project); err != nil {
		return err
	}
	return nil
}

package project

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entities"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/ports"
)

type Writer interface {
	Write(entities.Project) error
}

type FileWriter struct {
	repository ports.Repository
}

func NewWriter(repository ports.Repository) Writer {
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

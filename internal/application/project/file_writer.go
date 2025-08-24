package project

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/project/port"
)

type Writer interface {
	Write(entity.Project) error
}

type FileWriter struct {
	repository port.Repository
}

func NewWriter(repository port.Repository) Writer {
	return &FileWriter{
		repository: repository,
	}
}

func (cs *FileWriter) Write(project entity.Project) error {
	if err := cs.repository.Save(project); err != nil {
		return err
	}
	return nil
}

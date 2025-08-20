package configuration

import (
	"github.com/jairoprogramador/fastdeploy/internal/application/configuration/ports"
	"github.com/jairoprogramador/fastdeploy/internal/domain/configuration/entities"
	domain "github.com/jairoprogramador/fastdeploy/internal/domain/configuration/ports"
)

type FileWriter struct {
	repository domain.Repository
}

func NewWriter(repository domain.Repository) ports.Writer {
	return &FileWriter{
		repository: repository,
	}
}

func (cs *FileWriter) Write(config entities.Configuration) error {
	if err := cs.repository.Save(config); err != nil {
		return err
	}
	return nil
}

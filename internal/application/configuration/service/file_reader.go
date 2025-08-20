package service

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/configuration/entities"
	"github.com/jairoprogramador/fastdeploy/internal/domain/configuration/ports"
)

type Reader interface {
	Read() (entities.Configuration, error)
}

type FileReader struct {
	repository ports.Repository
}

func NewReader(repository ports.Repository) Reader {
	return &FileReader{
		repository: repository,
	}
}

func (cs *FileReader) Read() (entities.Configuration, error) {
	existsFile, err := cs.repository.Exists()
	if err != nil {
		return entities.Configuration{}, err
	}

	if !existsFile {
		return entities.NewDefaultConfiguration(), nil
	}

	config, err := cs.repository.Load()
	if err != nil {
		return entities.Configuration{}, err
	}

	return config, nil
}

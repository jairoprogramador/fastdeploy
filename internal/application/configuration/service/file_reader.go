package service

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/configuration/entity"
	"github.com/jairoprogramador/fastdeploy/internal/domain/configuration/port"
)

type Reader interface {
	Read() (entity.Configuration, error)
}

type FileReader struct {
	repository port.Repository
}

func NewReader(repository port.Repository) Reader {
	return &FileReader{
		repository: repository,
	}
}

func (cs *FileReader) Read() (entity.Configuration, error) {
	existsFile, err := cs.repository.Exists()
	if err != nil {
		return entity.Configuration{}, err
	}

	if !existsFile {
		return entity.NewDefaultConfiguration(), nil
	}

	config, err := cs.repository.Load()
	if err != nil {
		return entity.Configuration{}, err
	}

	return config, nil
}

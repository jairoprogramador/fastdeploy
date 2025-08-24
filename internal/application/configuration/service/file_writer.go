package service

import (
	"github.com/jairoprogramador/fastdeploy/internal/application/configuration/dto"
	"github.com/jairoprogramador/fastdeploy/internal/application/configuration/mapper"
	"github.com/jairoprogramador/fastdeploy/internal/domain/configuration/port"
)

type Writer interface {
	Write(dto.ConfigDto) error
}

type FileWriter struct {
	repository port.Repository
}

func NewWriter(repository port.Repository) Writer {
	return &FileWriter{
		repository: repository,
	}
}

func (cs *FileWriter) Write(dto dto.ConfigDto) error {
	config, err := mapper.ToDomain(dto)
	if err != nil {
		return err
	}

	if err := cs.repository.Save(config); err != nil {
		return err
	}
	return nil
}

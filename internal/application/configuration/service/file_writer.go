package service

import (
	"github.com/jairoprogramador/fastdeploy/internal/application/configuration/dto"
	"github.com/jairoprogramador/fastdeploy/internal/application/configuration/mapper"
	"github.com/jairoprogramador/fastdeploy/internal/domain/configuration/ports"
)

type Writer interface {
	Write(dto.ConfigDto) error
}

type FileWriter struct {
	repository ports.Repository
}

func NewWriter(repository ports.Repository) Writer {
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

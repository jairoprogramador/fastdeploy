package configuration

import (
	"github.com/jairoprogramador/fastdeploy/internal/application/configuration/ports"
	domainEntity "github.com/jairoprogramador/fastdeploy/internal/domain/configuration/entities"
	domainPort "github.com/jairoprogramador/fastdeploy/internal/domain/configuration/ports"
	domainService "github.com/jairoprogramador/fastdeploy/internal/domain/configuration/services"
)

type FileReader struct {
	repository domainPort.Repository
	validator domainService.Validator
}

func NewReader(repository domainPort.Repository, validator domainService.Validator) ports.Reader {
	return &FileReader{
		repository: repository,
		validator: validator,
	}
}

func (cs *FileReader) Read() (domainEntity.Configuration, error) {
	existsFile, err := cs.repository.Exists()
	if err != nil {
		return domainEntity.Configuration{}, err
	}

	if !existsFile {
		return domainEntity.NewDefaultConfiguration(), nil
	}

	config, err := cs.repository.Load()
	if err != nil {
		return domainEntity.Configuration{}, err
	}

	if err := cs.validator.Validate(config); err != nil {
		return domainEntity.Configuration{}, err
	}

	return config, nil
}

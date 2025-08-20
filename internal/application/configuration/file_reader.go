package configuration

import (
	"github.com/jairoprogramador/fastdeploy/internal/application/configuration/ports"
	domainEntity "github.com/jairoprogramador/fastdeploy/internal/domain/configuration/entities"
	domainPort "github.com/jairoprogramador/fastdeploy/internal/domain/configuration/ports"
)

type FileReader struct {
	repository domainPort.Repository
}

func NewReader(repository domainPort.Repository) ports.Reader {
	return &FileReader{
		repository: repository,
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

	return config, nil
}

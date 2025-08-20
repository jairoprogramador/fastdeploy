package project

import (
	domainEntity "github.com/jairoprogramador/fastdeploy/internal/domain/project/entities"
	domainPort "github.com/jairoprogramador/fastdeploy/internal/domain/project/ports"
	domainService "github.com/jairoprogramador/fastdeploy/internal/domain/project/services"
	"github.com/jairoprogramador/fastdeploy/internal/application/project/ports"
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

func (cs *FileReader) Read() (domainEntity.Project, error) {
	project, err := cs.repository.Load()
	if err != nil {
		return domainEntity.Project{}, err
	}

	if err := cs.validator.Validate(project); err != nil {
		return domainEntity.Project{}, err
	}

	return project, nil
}

func (cs *FileReader) ExistsFile() (bool, error) {
	return cs.repository.Exists()
}
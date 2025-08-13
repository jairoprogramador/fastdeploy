package project

import (
	"fmt"
)

type ProjectServiceImpl struct {
	repository ProjectRepository
	validator  ProjectValidator
}

func NewProjectService(
	repository ProjectRepository,
	validator ProjectValidator,
) ProjectService {
	return &ProjectServiceImpl{
		repository: repository,
		validator:  validator,
	}
}

func (ps *ProjectServiceImpl) Load() (*ProjectEntity, error) {
	project, err := ps.repository.Load()
	if err != nil {
		return nil, err
	}
	if err := ps.validate(*project); err != nil {
		return nil, fmt.Errorf("proyecto inválido: %w", err)
	}

	return project, nil
}

func (ps *ProjectServiceImpl) Save(projectEntity ProjectEntity) error {
	if err := ps.validate(projectEntity); err != nil {
		return fmt.Errorf("proyecto inválido: %w", err)
	}
	return ps.repository.Save(projectEntity)
}

func (ps *ProjectServiceImpl) Exists() bool {
	return ps.repository.Exists()
}

/*
func (ps *ProjectService) delete() error {
	return ps.repository.Delete()
} */

func (ps *ProjectServiceImpl) validate(project ProjectEntity) error {
	return ps.validator.Validate(project)
}

/* func Save(projectEntity ProjectEntity) error {
	return fmt.Errorf("función legacy - requiere implementación con factory")
}

func Load() (*ProjectEntity, error) {
	return nil, fmt.Errorf("función legacy - requiere implementación con factory")
}

func getProjectName() (string, error) {
	return "", fmt.Errorf("función legacy - requiere implementación con factory")
} */

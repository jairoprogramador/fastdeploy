package project

import (
	"fmt"

	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/project"
)

type ProjectCreatorImpl struct {
	pathResolver ProjectPathResolver
	idGenerator  ProjectIDGenerator
}

func NewProjectCreator(
	pathResolver ProjectPathResolver,
	idGenerator ProjectIDGenerator,
) domain.ProjectCreator {
	return &ProjectCreatorImpl{
		pathResolver: pathResolver,
		idGenerator:  idGenerator,
	}
}

func (pc *ProjectCreatorImpl) Create() (*domain.ProjectEntity, error) {
	projectName, err := pc.pathResolver.GetProjectName()
	if err != nil {
		return nil, fmt.Errorf("error al obtener nombre del proyecto: %w", err)
	}

	projectID, err := pc.idGenerator.GenerateID(projectName)
	if err != nil {
		return nil, fmt.Errorf("error al generar ID único: %w", err)
	}

	projectEntity := &domain.ProjectEntity{
		ProjectID:   projectID,
		ProjectName: projectName,
	}

	return projectEntity, nil
}

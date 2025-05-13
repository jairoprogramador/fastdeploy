// Package repository implementa el patrón Repository para el manejo de proyectos
package repository

import (
	"deploy/internal/domain/model"
	"deploy/internal/domain/repository"
	"deploy/internal/infrastructure/filesystem"
	"sync"
)

type projectRepositoryImpl struct {}

var (
	instanceProjectRepository     repository.ProjectRepository
	instanceOnceProjectRepository sync.Once
)

func GetProjectRepository() repository.ProjectRepository {
	instanceOnceProjectRepository.Do(func() {
		instanceProjectRepository = &projectRepositoryImpl {}
	})
	return instanceProjectRepository
}

func (st *projectRepositoryImpl) Load(path string) (model.Project, error) {
	return filesystem.LoadFromYAML[model.Project](path)
}

func (st *projectRepositoryImpl) Create(path string, project *model.Project) error {
	return filesystem.SaveToYAML(project, path)
}

func (s *projectRepositoryImpl) GetProjectName() (string, error) {
	return filesystem.GetParentDirectory()
}

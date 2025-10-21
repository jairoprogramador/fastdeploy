package application

import (
	"errors"
	"context"

	appPor "github.com/jairoprogramador/fastdeploy-core/internal/application/ports"
)

type RevisionProjectService struct {
	gitManager appPor.GitManager
}

func NewRevisionProjectService(gitManager appPor.GitManager) *RevisionProjectService {
	return &RevisionProjectService{
		gitManager: gitManager,
	}
}

func (s *RevisionProjectService) LoadProjectRevision(ctx context.Context, revisionDefault string, stepFinalName string) (string, error) {
	isGit, err := s.gitManager.IsGit()
	if err != nil {
		return "", err
	}

	if isGit {
		existChanges, err := s.gitManager.ExistChanges(ctx)
		if err != nil {
			return "", err
		}
		if existChanges {
			if stepFinalName == "test" {
				return revisionDefault, nil
			}

			return "", errors.New("hay cambios en el proyecto, ejecute 'git commit' primero")
		}
		return s.gitManager.GetCommitHash(ctx)
	}

	if stepFinalName != "test"{
		return "", errors.New("el projecto no esta configurado como repositorio git, ejecute 'git init' primero")
	}

	return revisionDefault, nil
}
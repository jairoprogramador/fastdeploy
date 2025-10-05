package application

import (
	"errors"
	"context"

	appports "github.com/jairoprogramador/fastdeploy/internal/application/ports"
)

type RevisionProjectService struct {
	gitManager appports.GitManager
}

func NewRevisionProjectService(gitManager appports.GitManager) *RevisionProjectService {
	return &RevisionProjectService{
		gitManager: gitManager,
	}
}

func (s *RevisionProjectService) LoadProjectRevision(ctx context.Context, revisionDefault string, stepFinalName string) (string, error) {
	git, err := s.gitManager.IsGit()
	if err != nil {
		return "", err
	}

	if git {
		existChanges, err := s.gitManager.ExistChanges(ctx)
		if err != nil {
			return "", err
		}
		if existChanges {
			if stepFinalName == "test"{
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
package application

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/jairoprogramador/fastdeploy-core/internal/application/ports"
)

type LocalGitService struct {
	executor ports.CommandService
}

func NewLocalGitService(
	executor ports.CommandService) ports.GitService {

	return &LocalGitService{
		executor: executor,
	}
}

func (g *LocalGitService) IsGit(pathProject string) (bool, error) {
	_, err := os.Stat(filepath.Join(pathProject, ".git"))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (g *LocalGitService) GetCommitHash(ctx context.Context, pathProject string) (string, error) {
	hash, _, err := g.executor.Run(ctx, pathProject, "git rev-parse --short HEAD")
	if err != nil {
		return "", err
	}
	hash = strings.TrimSpace(hash)
	hash = strings.ReplaceAll(hash, "|", "")
	hash = strings.ReplaceAll(hash, "\n", "")
	return hash, nil
}

func (g *LocalGitService) ExistChanges(ctx context.Context, pathProject string) (bool, error) {

	_, codeExit, err := g.executor.Run(ctx, pathProject, "git diff --quiet")
	if err != nil {
		return false, err
	}
	if codeExit != 0 {
		return true, nil
	}

	_, codeExit, err = g.executor.Run(ctx, pathProject, "git diff --cached --quiet")
	if err != nil {
		return false, err
	}
	if codeExit != 0 {
		return true, nil
	}

	log, codeExit, err := g.executor.Run(ctx, pathProject, "git ls-files --others --exclude-standard")
	if err != nil {
		return false, err
	}
	if codeExit != 0 {
		return false, nil
	}
	if len(log) > 0 {
		return true, nil
	}

	return false, nil
}

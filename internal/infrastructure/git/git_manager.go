package git

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"github.com/jairoprogramador/fastdeploy/internal/application/ports"
)

type GitManager struct{
	executor ports.CommandExecutor
	pathAppProject string
}

func NewGitManager(executor ports.CommandExecutor, pathAppProject string) ports.GitManager {
	return &GitManager{
		executor: executor,
		pathAppProject: pathAppProject,
	}
}

func (g *GitManager) IsGit() (bool, error) {
	_, err := os.Stat(filepath.Join(g.pathAppProject, ".git"))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (g *GitManager) GetCommitHash(ctx context.Context) (string, error) {
	hash, _, err := g.executor.Execute(ctx, "", "git rev-parse --short HEAD")
	if err != nil {
		return "", err
	}
	hash = strings.TrimSpace(hash)
	hash = strings.ReplaceAll(hash, "|", "")
	hash = strings.ReplaceAll(hash, "\n", "")
	return hash, nil
}

func (g *GitManager) ExistChanges(ctx context.Context) (bool, error) {
	_, codeExit, err := g.executor.Execute(ctx, g.pathAppProject, "git diff --quiet")
	if err != nil {
		return false, err
	}
	return codeExit == 1, nil
}
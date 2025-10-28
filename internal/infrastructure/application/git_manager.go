package application

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"github.com/jairoprogramador/fastdeploy-core/internal/application/ports"
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
	if codeExit != 0 {
		return true, nil
	}

	_, codeExit, err = g.executor.Execute(ctx, g.pathAppProject, "git diff --cached --quiet")
	if err != nil {
		return false, err
	}
	if codeExit != 0 {
		return true, nil
	}

	log, codeExit, err := g.executor.Execute(ctx, g.pathAppProject, "git ls-files --others --exclude-standard")
	if err != nil {
		return false, err
	}
	if len(log) > 0 {
		return true, nil
	}

	return false, nil
}
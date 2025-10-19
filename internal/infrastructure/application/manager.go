package application

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"github.com/jairoprogramador/fastdeploy/internal/application/ports"
)

type Manager struct {
	pathProjectRootFastDeploy string
	pathRepositoryRootFastDeploy string
	projectName string
	repositoryName string
	environment string
	pathRepository string
	pathProject string
}

func NewManager(
	pathProjectRootFastDeploy,
	pathRepositoryRootFastDeploy,
	projectName,
	repositoryName,
	environment string) (ports.WorkspaceManager, error) {

	if pathProjectRootFastDeploy == "" {
		return nil, fmt.Errorf("path project root fast deploy is required")
	}
	if pathRepositoryRootFastDeploy == "" {
		return nil, fmt.Errorf("path repository root fast deploy is required")
	}
	if projectName == "" {
		return nil, fmt.Errorf("project name is required")
	}
	if repositoryName == "" {
		return nil, fmt.Errorf("repository name is required")
	}
	if environment == "" {
		return nil, fmt.Errorf("environment is required")
	}

	pathRepository := filepath.Join(pathRepositoryRootFastDeploy, repositoryName)
	pathProject := filepath.Join(pathProjectRootFastDeploy, projectName)

	return &Manager {
		pathProjectRootFastDeploy: pathProjectRootFastDeploy,
		pathRepositoryRootFastDeploy: pathRepositoryRootFastDeploy,
		projectName: projectName,
		repositoryName: repositoryName,
		environment: environment,
		pathRepository: pathRepository,
		pathProject: pathProject,
	}, nil
}

func (m *Manager) Prepare(stepName string) (string, error) {
	pathStepSource, err := m.getPathStepSource(stepName)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	} else if os.IsNotExist(err) {
		return "", nil
	}

	destPath := filepath.Join(m.pathProject, m.repositoryName, m.environment, stepName)

	if err := os.MkdirAll(destPath, 0755); err != nil {
		return "", fmt.Errorf("error al crear el workspace del paso: %w", err)
	}

	if err := copyDir(pathStepSource, destPath); err != nil {
		return "", fmt.Errorf("error al copiar los archivos de la plantilla al workspace: %w", err)
	}

	return destPath, nil
}

func (m *Manager) getPathStepSource(stepName string) (string, error) {
	pathSteps := filepath.Join(m.pathRepository, "steps")
	entries, err := os.ReadDir(pathSteps)
	if err != nil {
		return "", err
	}

	regex := regexp.MustCompile(fmt.Sprintf(`^\d+-%s$`, regexp.QuoteMeta(stepName)))
	for _, entry := range entries {
		if entry.IsDir() && regex.MatchString(entry.Name()) {
			return filepath.Join(pathSteps, entry.Name()), nil
		}
	}

	return "", os.ErrNotExist
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		return copyFile(path, dstPath, info)
	})
}

func copyFile(src, dst string, info os.FileInfo) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

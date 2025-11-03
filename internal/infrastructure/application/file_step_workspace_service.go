package application

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/jairoprogramador/fastdeploy-core/internal/application/dto"
	"github.com/jairoprogramador/fastdeploy-core/internal/application/ports"
)

type FileStepWorkspaceService struct {
	pathProjectRoot    string
	pathRepositoryRoot string
}

func NewFileStepWorkspaceService(
	pathProjectRoot,
	pathRepositoryRoot string,
) ports.StepWorkspaceService {
	return &FileStepWorkspaceService{
		pathProjectRoot:    pathProjectRoot,
		pathRepositoryRoot: pathRepositoryRoot,
	}
}

func (m *FileStepWorkspaceService) Prepare(namesRequest dto.NamesParams, runRequest dto.RunParams) (string, error) {
	pathStepRepository, err := m.getPathStepRepository(runRequest.StepName(), namesRequest.RepositoryName())
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	pathStepProject := m.getPathStepProject(namesRequest, runRequest)

	if err := os.MkdirAll(pathStepProject, 0755); err != nil {
		return "", fmt.Errorf("error al crear el workspace del paso: %w", err)
	}

	if err := copyDir(pathStepRepository, pathStepProject); err != nil {
		return "", fmt.Errorf("error al copiar los archivos de la plantilla al workspace: %w", err)
	}

	return pathStepProject, nil
}

func (m *FileStepWorkspaceService) getPathStepProject(namesRequest dto.NamesParams, runRequest dto.RunParams) string {
	pathProject := filepath.Join(m.pathProjectRoot, namesRequest.ProjectName())
	return filepath.Join(pathProject, namesRequest.RepositoryName(), runRequest.Environment(), runRequest.StepName())
}

func (m *FileStepWorkspaceService) getPathStepRepository(stepName string, repositoryName string) (string, error) {
	pathRepository := filepath.Join(m.pathRepositoryRoot, repositoryName)

	pathSteps := filepath.Join(pathRepository, "steps")
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

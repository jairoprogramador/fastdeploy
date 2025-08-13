package project

import (
	"fmt"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy/internal/adapters/filesystem"
	"github.com/jairoprogramador/fastdeploy/internal/constants"
)

type ProjectPathResolver interface {
	GetProjectName() (string, error)
	GetProjectPath() (string, error)
}

type ProjectPathResolverImpl struct {
	workingDir filesystem.WorkingDirectory
}

func NewProjectPathResolver(workingDir filesystem.WorkingDirectory) ProjectPathResolver {
	return &ProjectPathResolverImpl{
		workingDir: workingDir,
	}
}

func (opr *ProjectPathResolverImpl) GetProjectName() (string, error) {
	dir, err := opr.workingDir.Getwd()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio de trabajo: %w", err)
	}
	return filepath.Base(dir), nil
}

func (opr *ProjectPathResolverImpl) GetProjectPath() (string, error) {
	dir, err := opr.workingDir.Getwd()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio de trabajo: %w", err)
	}
	return filepath.Join(dir, constants.ProjectFileName), nil
}

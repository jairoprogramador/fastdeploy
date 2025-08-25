package service

import (
	"os"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy/internal/domain/project/port"
)

type ProjectName struct{}

func NewProjectName() port.Name {
	return &ProjectName{}
}

func (s *ProjectName) GetName() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Base(dir), nil
}

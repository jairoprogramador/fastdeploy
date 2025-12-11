package vos

import (
	"fmt"
	"os"
	"path/filepath"
)

const defaultRootDir = ".fastdeploy"

type RootPath struct {
	path string
}

func NewRootPath(path string) (RootPath, error) {
	if path == "" {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return RootPath{}, fmt.Errorf("could not get user home directory: %w", err)
		}
		path = filepath.Join(userHomeDir, defaultRootDir)
	}
	return RootPath{path: path}, nil
}

func (r RootPath) Path() string {
	return r.path
}

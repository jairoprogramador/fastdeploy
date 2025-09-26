package command

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jairoprogramador/fastdeploy/internal/domain/command/port"
)

type WorkdirFile struct{}

func NewWorkdirFile() port.WorkdirPort {
	return &WorkdirFile{}
}

func (w *WorkdirFile) Copy(sourcePath string, destinationPath string) error {
	sourceInfo, err := os.Stat(sourcePath)
	if err != nil {
		return err
	}

	if sourceInfo.IsDir() {
		return w.copyDirectoryRecursive(sourcePath, destinationPath)
	}

	destInfo, err := os.Stat(destinationPath)
	if err == nil && destInfo.IsDir() {
		destinationPath = filepath.Join(destinationPath, filepath.Base(sourcePath))
	}

	return w.copyFile(sourcePath, destinationPath)
}

func (w *WorkdirFile) copyFile(sourceFilePath string, destinationFilePath string) error {
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	if err := os.MkdirAll(filepath.Dir(destinationFilePath), os.ModePerm); err != nil {
		return err
	}

	destFile, err := os.Create(destinationFilePath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

func (w *WorkdirFile) copyDirectoryRecursive(sourcePath string, destinationPath string) error {
	const maxDepth = 3

	return filepath.Walk(sourcePath, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourcePath, currentPath)
		if err != nil {
			return err
		}

		if relPath == "." {
			return nil
		}

		depth := strings.Count(relPath, string(os.PathSeparator))

		if info.IsDir() && depth >= maxDepth {
			return filepath.SkipDir
		}

		destPath := filepath.Join(destinationPath, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		return w.copyFile(currentPath, destPath)
	})
}

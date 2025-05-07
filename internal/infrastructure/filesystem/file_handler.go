package filesystem

import (
	"os"
	"path/filepath"
)

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}

func RemoveFile(filePath string) error {
	if err := os.Remove(filePath); !os.IsNotExist(err) {
		return err
	}
	return nil
}

func CreateFile(filePath string) (*os.File, error) {
	return os.Create(filePath)
}

func WriteFile(filePath, content string) error {
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return err
	}
	return nil
}

func OpenFile(filePath string) (*os.File, error) {
	return os.Open(filePath)
}

func DirectoryExists(filePath string) (bool, error) {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}

func CreateDirectory(nameDirectory string) error {
	if err := os.MkdirAll(nameDirectory, 0755); err != nil {
		return err
	}
	return nil
}

func RecreateDirectory(nameDirectory string) error {
	if _, err := os.Stat(nameDirectory); err == nil {
		if err := os.RemoveAll(nameDirectory); err != nil {
			return err
		}
	}
	return CreateDirectory(nameDirectory)
}

func CompletePermits(nameDirectory string) error {
	if err := os.Chmod(nameDirectory, 0777); err != nil {
		return err
	}
	return nil
}

func CreateDirectoryFilePath(filePath string) error {
	exists, err := DirectoryExists(filePath)
	if !exists {
		directory := GetDirectory(filePath)
		return CreateDirectory(directory)
	}
	return err
}

func GetDirectory(filePath string) string {
	return filepath.Dir(filePath)
}

func GetParentDirectory() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	projectId := filepath.Base(currentDir)
	return projectId, nil
}

func GetProjectDirectory() (string, error) {
	return os.Getwd()
}

func GetHomeDirectory() (string, error) {
	return os.UserHomeDir()
}

func GetPath(paths ...string) string {
	return filepath.Join(paths...)
}

package filesystem

import (
	"os"
	"path/filepath"
)

func ExistsFile(filePath string) (bool, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !fileInfo.IsDir(), nil
}

func RemoveFile(filePath string) error {
	exists, err := ExistsFile(filePath)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	if err := os.Remove(filePath); err != nil {
		return err
	}
	return nil
}

func CreateFile(filePath string) (*os.File, error) {
	dir := GetDirectory(filePath)
	exists, _ := ExistsDirectory(dir)
	if !exists {
		if err := CreateDirectory(dir); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
}

func WriteFile(filePath, content string) error {
	file, err := CreateFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := file.WriteString(content); err != nil {
		return err
	}
	return nil
}

func OpenFile(filePath string) (*os.File, error) {
	return os.Open(filePath)
}

func ExistsDirectory(dirPath string) (bool, error) {
	fileInfo, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func CreateDirectory(dirPath string) error {
	if err := os.MkdirAll(dirPath, 0755); err != nil {
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

func GetDirectory(pathFile string) string {
	return filepath.Dir(pathFile)
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

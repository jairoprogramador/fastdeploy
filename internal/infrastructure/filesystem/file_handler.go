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

func CreateFile(pathFile string) (*os.File, error) {
	dir := filepath.Dir(pathFile)
	exists, _ := ExistsDirectory(dir)
	if !exists {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(pathFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
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

func GetPath(paths ...string) string {
	return filepath.Join(paths...)
}

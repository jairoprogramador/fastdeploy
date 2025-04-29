package filesystem

import "path/filepath"
import "os"

func FileExists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}
	return true
}
func Removefile(filePath string) error { 
	err := os.Remove(filePath)
	return err
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

func GetParentDirectory() (string, error){
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	projectId := filepath.Base(currentDir)
	return projectId, nil
}

func GetHomeDirectory() (string, error){
	return os.UserHomeDir()
}

func GetPath(directory, file string) string {
	return filepath.Join(directory, file)
}
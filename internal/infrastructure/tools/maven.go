package tools

import (
	"deploy/internal/domain/constant"
	"deploy/internal/infrastructure/filesystem"
	"fmt"
	"os"
	"strings"
)

func GetMavenVersion() (string, error) {
	return ExecuteCommand("mvn", "--version")
}

func CleanAndPackage() (string, error) {
	return ExecuteCommand("mvn", "clean", "package")
}

func GetFullPathFiles(directory string) ([]string, error) {
	var pathFiles []string

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, archivo := range files {
		if !archivo.IsDir() && strings.HasSuffix(archivo.Name(), ".jar") &&
			!strings.Contains(archivo.Name(), "sources") &&
			!strings.Contains(archivo.Name(), "original") {
			path := filesystem.GetPath(directory,archivo.Name())	
			pathFiles = append(pathFiles, path)
		}
	}

	if len(pathFiles) <= 0 {
		return nil, fmt.Errorf(constant.MessageErrorNoPackedFile)
	}

	return pathFiles, nil
}
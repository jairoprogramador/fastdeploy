package tools

import (
	"os"
	"strings"
	"fmt"
	"deploy/internal/domain"
)

func GetMavenVersion() (string, error) {
	return ExecuteCommand("mvn", "--version")
}

func CleanAndPackage() (string, error) {
	return ExecuteCommand("mvn", "clean", "package")
}

func SearchFile(directory string) ([]string, error) {
	var jarFiles []string

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, archivo := range files {
		if !archivo.IsDir() && strings.HasSuffix(archivo.Name(), ".jar") &&
			!strings.Contains(archivo.Name(), "sources") &&
			!strings.Contains(archivo.Name(), "original") {
			jarFiles = append(jarFiles, directory+"/"+archivo.Name())
		}
	}

	if len(jarFiles) <= 0 {
		return nil, fmt.Errorf(constants.MessageErrorNoPackedFile)
	}

	return jarFiles, nil
}
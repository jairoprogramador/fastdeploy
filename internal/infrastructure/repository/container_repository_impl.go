package repository

import (
	"deploy/internal/domain/repository"
	"deploy/internal/infrastructure/filesystem"
	"text/template"
	"fmt"
	"strings"
	"os"
	"sync"
)

type containerRepositoryImpl struct {
	fileRepository repository.FileRepository
}

var (
	instanceContainerRepository     repository.ContainerRepository
	instanceOnceContainerRepository sync.Once
)

func GetContainerRepository() repository.ContainerRepository {
	instanceOnceContainerRepository.Do(func() {
		instanceContainerRepository = &containerRepositoryImpl {
			fileRepository: GetFileRepository(),
		}
	})
	return instanceContainerRepository
}

func (st *containerRepositoryImpl) GetFullPathResource() (string, error) {
	directoryTarget := "target"
	exists := st.fileRepository.ExistsDirectory(directoryTarget)
	if !exists {
		return "", fmt.Errorf("no se encontró el directorio target")
	}

	fullPathJarFiles, err := getFullPathResources(directoryTarget)
	if err != nil {
		return "", err
	}
	return fullPathJarFiles[0], nil
}

func (st *containerRepositoryImpl) GetContentTemplate(pathTemplate string, params any) (string, error) {
	dockerfileTemplate, err := template.ParseFiles(pathTemplate)
	if err != nil {
		return "", err
	}
	var result strings.Builder
	err = dockerfileTemplate.Execute(&result, params)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

func getFullPathResources(directory string) ([]string, error) {
	var pathFiles []string

	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	for _, archivo := range files {
		if !archivo.IsDir() && strings.HasSuffix(archivo.Name(), ".jar") &&
			!strings.Contains(archivo.Name(), "sources") &&
			!strings.Contains(archivo.Name(), "original") {
				
			path := filesystem.GetPath(directory, archivo.Name())
			path = filesystem.GetAbsolutePath(path)
			pathFiles = append(pathFiles, path)
		}
	}

	if len(pathFiles) <= 0 {
		return nil, fmt.Errorf("no se encontró el archivo jar")
	}

	return pathFiles, nil
}

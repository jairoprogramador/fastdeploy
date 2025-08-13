package project

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/filesystem"
	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/project"
	"path/filepath"
)

type YAMLProjectRepositoryImpl struct {
	fileSystem   filesystem.FileSystem
	pathResolver ProjectPathResolver
	serializer   YAMLProjectSerializer
}

func NewYAMLProjectRepository(
	fileSystem filesystem.FileSystem,
	pathResolver ProjectPathResolver,
	serializer YAMLProjectSerializer,
) domain.ProjectRepository {
	return &YAMLProjectRepositoryImpl{
		fileSystem:   fileSystem,
		pathResolver: pathResolver,
		serializer:   serializer,
	}
}

func (ypr *YAMLProjectRepositoryImpl) Save(projectEntity domain.ProjectEntity) error {
	filePath, err := ypr.pathResolver.GetProjectPath()
	if err != nil {
		return fmt.Errorf("error al obtener ruta del archivo: %w", err)
	}

	dir := filepath.Dir(filePath)
	if err := ypr.fileSystem.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error al crear el directorio del proyecto: %w", err)
	}

	data, err := ypr.serializer.Serialize(projectEntity)
	if err != nil {
		return fmt.Errorf("error al serializar el proyecto: %w", err)
	}

	if err := ypr.fileSystem.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error al guardar el proyecto: %w", err)
	}

	return nil
}

func (ypr *YAMLProjectRepositoryImpl) Load() (*domain.ProjectEntity, error) {
	filePath, err := ypr.pathResolver.GetProjectPath()
	if err != nil {
		return nil, fmt.Errorf("error al obtener ruta del archivo: %w", err)
	}

	data, err := ypr.fileSystem.ReadFile(filePath)
	if err != nil {
		if ypr.fileSystem.IsNotExist(err) {
			return &domain.ProjectEntity{}, nil
		}
		return nil, fmt.Errorf("error al leer el archivo del proyecto: %w", err)
	}

	project, err := ypr.serializer.Deserialize(data)
	if err != nil {
		return nil, fmt.Errorf("error al deserializar el proyecto: %w", err)
	}

	return project, nil
}

func (ypr *YAMLProjectRepositoryImpl) Exists() bool {
	filePath, err := ypr.pathResolver.GetProjectPath()
	if err != nil {
		return false
	}

	_, err = ypr.fileSystem.ReadFile(filePath)
	return err == nil
}

func (ypr *YAMLProjectRepositoryImpl) Delete() error {
	_, err := ypr.pathResolver.GetProjectPath()
	if err != nil {
		return fmt.Errorf("error al obtener ruta del archivo: %w", err)
	}

	return fmt.Errorf("método Delete no implementado - requiere implementación en filesystem.FileSystem")
}

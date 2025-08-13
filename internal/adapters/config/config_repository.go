package config

import (
	"fmt"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy/internal/adapters/filesystem"
	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/config"
)

type YAMLConfigRepositoryImpl struct {
	fileSystem   filesystem.FileSystem
	userSystem   filesystem.UserSystem
	pathResolver ConfigPathResolver
	serializer   YAMLConfigSerializer
}

func NewYAMLConfigRepository(
	fileSystem filesystem.FileSystem,
	userSystem filesystem.UserSystem,
	pathResolver ConfigPathResolver,
	serializer YAMLConfigSerializer,
) domain.ConfigRepository {
	return &YAMLConfigRepositoryImpl{
		fileSystem:   fileSystem,
		userSystem:   userSystem,
		pathResolver: pathResolver,
		serializer:   serializer,
	}
}

func (cr *YAMLConfigRepositoryImpl) Save(configEntity domain.ConfigEntity) error {
	filePath, err := cr.pathResolver.GetConfigFilePath()
	if err != nil {
		return fmt.Errorf("error al obtener ruta del archivo: %w", err)
	}

	dir := filepath.Dir(filePath)
	if err := cr.fileSystem.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error al crear el directorio de configuración: %w", err)
	}

	data, err := cr.serializer.Serialize(configEntity)
	if err != nil {
		return fmt.Errorf("error al serializar la configuración: %w", err)
	}

	if err := cr.fileSystem.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error al guardar la configuración: %w", err)
	}

	return nil
}

func (cr *YAMLConfigRepositoryImpl) Load() (*domain.ConfigEntity, error) {
	filePath, err := cr.pathResolver.GetConfigFilePath()
	if err != nil {
		return nil, fmt.Errorf("error al obtener ruta del archivo: %w", err)
	}

	data, err := cr.fileSystem.ReadFile(filePath)
	if err != nil {
		if cr.fileSystem.IsNotExist(err) {
			return &domain.ConfigEntity{}, nil
		}
		return nil, fmt.Errorf("error al leer el archivo de configuración: %w", err)
	}

	config, err := cr.serializer.Deserialize(data)
	if err != nil {
		return nil, fmt.Errorf("error al deserializar la configuración: %w", err)
	}

	return config, nil
}

func (cr *YAMLConfigRepositoryImpl) Exists() bool {
	filePath, err := cr.pathResolver.GetConfigFilePath()
	if err != nil {
		return false
	}

	_, err = cr.fileSystem.ReadFile(filePath)
	return err == nil
}

func (cr *YAMLConfigRepositoryImpl) Delete() error {
	_, err := cr.pathResolver.GetConfigFilePath()
	if err != nil {
		return fmt.Errorf("error al obtener ruta del archivo: %w", err)
	}

	// Asumiendo que el método correcto es RemoveFile o similar
	// Si no existe, podemos implementar la lógica directamente
	return fmt.Errorf("método Delete no implementado - requiere implementación en filesystem.FileSystem")
}

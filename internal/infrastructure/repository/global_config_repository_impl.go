package repository

import (
	"deploy/internal/domain/model"
	"deploy/internal/domain/repository"
	"deploy/internal/infrastructure/filesystem"
	"sync"
)

// globalConfigRepositoryImpl implementa la interfaz GlobalConfigRepository
// utilizando el patrón Singleton para asegurar una única instancia.
type globalConfigRepositoryImpl struct{}

var (
	instanceGlobalConfigRepository     repository.GlobalConfigRepository
	instanceOnceGlobalConfigRepository sync.Once
)

// NewGlobalConfigRepository crea una nueva instancia del repositorio de configuración global.
// Implementa el patrón Singleton para asegurar una única instancia.
func GetGlobalConfigRepository() repository.GlobalConfigRepository {
	instanceOnceGlobalConfigRepository.Do(func() {
		instanceGlobalConfigRepository = &globalConfigRepositoryImpl{}
	})
	return instanceGlobalConfigRepository
}

// Load carga la configuración global desde el archivo YAML.
func (st *globalConfigRepositoryImpl) Load() (model.GlobalConfig, error) {
	filePath, err := st.getPathGlobalConfigFile()
	if err != nil {
		return model.GlobalConfig{}, err
	}
	return filesystem.LoadFromYAML[model.GlobalConfig](filePath)
}

// ExistsFile verifica si existe el archivo de configuración global.
func (st *globalConfigRepositoryImpl) ExistsFile() bool {
	filePath, err := st.getPathGlobalConfigFile()
	if err != nil {
		return false
	}
	exists, _ := filesystem.ExistsFile(filePath)
	return exists
}

// RemoveFile elimina el archivo de configuración global.
func (st *globalConfigRepositoryImpl) RemoveFile() error {
	filePath, err := st.getPathGlobalConfigFile()
	if err != nil {
		return err
	}
	return filesystem.RemoveFile(filePath)
}

// Create crea un nuevo archivo de configuración global.
func (st *globalConfigRepositoryImpl) Create(globalConfig *model.GlobalConfig) error {
	if err := st.RemoveFile(); err != nil {
		return err
	}

	filePath, err := st.getPathGlobalConfigFile()
	if err != nil {
		return err
	}

	if err := filesystem.SaveToYAML(globalConfig, filePath); err != nil {
		return err
	}
	return nil
}

// getPathGlobalConfigFile obtiene la ruta completa del archivo de configuración global.
func (st *globalConfigRepositoryImpl) getPathGlobalConfigFile() (string, error) {
	homeDir, err := filesystem.GetHomeDirectory()
	if err != nil {
		return "", err
	}
	return filesystem.GetPath(homeDir, UserDirectory, GlobalConfigFile), nil
}

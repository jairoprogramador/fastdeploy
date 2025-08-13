package config

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/filesystem"
	"github.com/jairoprogramador/fastdeploy/internal/constants"
	"path/filepath"
)

type ConfigPathResolver interface {
	GetConfigFilePath() (string, error)
	GetConfigDirPath() (string, error)
}

type ConfigPathResolverImpl struct {
	userSystem filesystem.UserSystem
}

func NewConfigPathResolver(userSystem filesystem.UserSystem) ConfigPathResolver {
	return &ConfigPathResolverImpl{
		userSystem: userSystem,
	}
}

func (cpr *ConfigPathResolverImpl) GetConfigFilePath() (string, error) {
	configDirPath, err := cpr.GetConfigDirPath()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDirPath, constants.GlobalConfigFileName), nil
}

func (cpr *ConfigPathResolverImpl) GetConfigDirPath() (string, error) {
	currentUser, err := cpr.userSystem.Current()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio del usuario: %w", err)
	}
	return filepath.Join(currentUser.HomeDir, constants.FastDeployDirName), nil
}

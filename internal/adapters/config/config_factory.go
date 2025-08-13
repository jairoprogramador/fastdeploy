package config

import (
	"github.com/jairoprogramador/fastdeploy/internal/adapters/filesystem"
	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/config"
)

type ConfigFactory struct{}

func NewConfigFactory() *ConfigFactory {
	return &ConfigFactory{}
}

func (cf *ConfigFactory) CreateService() domain.ConfigService {
	userSystem := &filesystem.OSUserSystem{}
	fileSystem := &filesystem.OSFileSystem{}

	pathResolver := cf.CreatePathResolver()
	serializer := NewYAMLConfigSerializer()
	repository := NewYAMLConfigRepository(fileSystem, userSystem, pathResolver, serializer)
	validator := domain.NewConfigValidatorImpl()

	return domain.NewConfigService(repository, validator)
}

func (cf *ConfigFactory) CreatePathResolver() ConfigPathResolver {
	userSystem := &filesystem.OSUserSystem{}
	return NewConfigPathResolver(userSystem)
}

/* func (cf *ConfigFactory) CreateConfigService(
	fileSystem filesystem.FileSystem,
	userSystem filesystem.UserSystem,
) domain.ConfigService {
	pathResolver := NewConfigPathResolver(userSystem)
	serializer := NewYAMLConfigSerializer()
	repository := NewYAMLConfigRepository(fileSystem, userSystem, pathResolver, serializer)
	validator := domain.NewConfigValidatorImpl()

	return domain.NewConfigService(repository, validator)
} */

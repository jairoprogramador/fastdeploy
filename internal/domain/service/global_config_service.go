package service

import (
    "deploy/internal/domain/constant"
	"deploy/internal/domain/model"
	"deploy/internal/domain/repository"
	"deploy/internal/domain/router"
	"errors"
)

var (
	ErrGlobalConfigNotFound     = errors.New(constant.MsgGlobalConfigNotFound)
    ErrGlobalConfigCanNotBeNull = errors.New(constant.MsgGlobalConfigCannoBeNull)
)

type GlobalConfigServiceInterface interface {
	Load() (model.GlobalConfig, error)
	Create(globalConfig *model.GlobalConfig) error
}

type globalConfigService struct {
	yamlRepository           repository.YamlRepository
	fileRepository           repository.FileRepository
	router                   *router.Router
}

func NewGlobalConfigService(
	yamlRepository repository.YamlRepository,
	fileRepository repository.FileRepository,
	router *router.Router,
) GlobalConfigServiceInterface {
	return &globalConfigService{
			yamlRepository: yamlRepository,
			fileRepository: fileRepository,
		router:         router,
		}
}

func (s *globalConfigService) Load() (model.GlobalConfig, error) {
	path := s.router.GetFullPathGlobalConfigFile()

	if exists := s.fileRepository.ExistsFile(path); !exists {
		return model.GlobalConfig{}, ErrGlobalConfigNotFound
	}

	var globalConfig model.GlobalConfig
	err := s.yamlRepository.Load(path, &globalConfig)
	if err != nil {
		return model.GlobalConfig{}, err
	}

	return globalConfig, nil
}

func (s *globalConfigService) Create(globalConfig *model.GlobalConfig) error {
	if globalConfig == nil {
		return ErrGlobalConfigCanNotBeNull
	}

	path := s.router.GetFullPathGlobalConfigFile()

	if err := s.fileRepository.DeleteFile(path); err != nil {
		return err
	}

	if err := s.yamlRepository.Save(path, globalConfig); err != nil {
		return err
	}

	return nil
}

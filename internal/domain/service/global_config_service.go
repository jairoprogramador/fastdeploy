package service

import (
    "deploy/internal/domain/constant"
	"deploy/internal/domain/model"
	"deploy/internal/domain/repository"
	"errors"
	"sync"
)

// Errores personalizados del servicio de proyecto
var (
	ErrGlobalConfigNotFound    = errors.New(constant.MsgGlobalConfigNotFound)
    ErrGlobalConfigCannoBeNull = errors.New(constant.MsgGlobalConfigCannoBeNull)
)

// GlobalConfigServiceInterface define la interfaz para el servicio de configuración global
type GlobalConfigServiceInterface interface {
	Load() (model.GlobalConfig, error)
	Create(globalConfig *model.GlobalConfig) (model.GlobalConfig, error)
	SetGlobalConfigRepository(globalConfigRepo repository.GlobalConfigRepository)
}

// GlobalConfigService maneja la lógica de negocio para la configuración global
type GlobalConfigService struct {
	globalConfigRepo 			repository.GlobalConfigRepository
    mutexGlobalConfigService 	sync.RWMutex
}

var (
	instanceGlobalConfigService     *GlobalConfigService
	instanceOnceGlobalConfigService sync.Once
)

// NewGlobalConfigService crea una nueva instancia del servicio de configuración global
func GetGlobalConfigService(globalConfigRepo repository.GlobalConfigRepository) GlobalConfigServiceInterface {
	instanceOnceGlobalConfigService.Do(func() {
		instanceGlobalConfigService = &GlobalConfigService{
			globalConfigRepo: globalConfigRepo,
		}
	})
	return instanceGlobalConfigService
}

// SetGlobalConfigRepository establece el repositorio de configuración global
func (s *GlobalConfigService) SetGlobalConfigRepository(globalConfigRepo repository.GlobalConfigRepository) {
	s.mutexGlobalConfigService.Lock()
	defer s.mutexGlobalConfigService.Unlock()
	s.globalConfigRepo = globalConfigRepo
}

// Load carga la configuración global desde el repositorio
func (s *GlobalConfigService) Load() (model.GlobalConfig, error) {
	if exists := s.globalConfigRepo.ExistsFile(); !exists {
		return model.GlobalConfig{}, ErrGlobalConfigNotFound
	}

	globalConfig, err := s.globalConfigRepo.Load()
	if err != nil {
		return model.GlobalConfig{}, ErrGlobalConfigNotFound
	}

	return globalConfig, nil
}

// Create crea una nueva configuración global
func (s *GlobalConfigService) Create(globalConfig *model.GlobalConfig) (model.GlobalConfig, error) {
	if globalConfig == nil {
		return model.GlobalConfig{}, ErrGlobalConfigCannoBeNull
	}

	if err := s.globalConfigRepo.Create(globalConfig); err != nil {
		return model.GlobalConfig{}, err
	}

	return *globalConfig, nil
}

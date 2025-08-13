package config

import (
	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/config"
	"gopkg.in/yaml.v3"
)

type YAMLConfigSerializer interface {
	Serialize(config domain.ConfigEntity) ([]byte, error)
	Deserialize(data []byte) (*domain.ConfigEntity, error)
}

type YAMLConfigSerializerImpl struct{}

func NewYAMLConfigSerializer() YAMLConfigSerializer {
	return &YAMLConfigSerializerImpl{}
}

func (ycs *YAMLConfigSerializerImpl) Serialize(config domain.ConfigEntity) ([]byte, error) {
	return yaml.Marshal(config)
}

func (ycs *YAMLConfigSerializerImpl) Deserialize(data []byte) (*domain.ConfigEntity, error) {
	var cfg domain.ConfigEntity
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

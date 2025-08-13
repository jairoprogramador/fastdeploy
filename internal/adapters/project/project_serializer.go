package project

import (
	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/project"
	"gopkg.in/yaml.v3"
)

type YAMLProjectSerializer interface {
	Serialize(project domain.ProjectEntity) ([]byte, error)
	Deserialize(data []byte) (*domain.ProjectEntity, error)
}

type YAMLProjectSerializerImpl struct{}

func NewYAMLProjectSerializer() YAMLProjectSerializer {
	return &YAMLProjectSerializerImpl{}
}

func (yps *YAMLProjectSerializerImpl) Serialize(project domain.ProjectEntity) ([]byte, error) {
	return yaml.Marshal(project)
}

func (yps *YAMLProjectSerializerImpl) Deserialize(data []byte) (*domain.ProjectEntity, error) {
	var project domain.ProjectEntity
	if err := yaml.Unmarshal(data, &project); err != nil {
		return nil, err
	}
	return &project, nil
}

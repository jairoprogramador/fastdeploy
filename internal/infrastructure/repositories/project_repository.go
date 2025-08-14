package repositories

import (
	"os"

	"gopkg.in/yaml.v3"
)

const PROJECT_FILE_NAME = "deploy.yaml"

type ProjectRepository struct{}

func NewProjectRepository() *ProjectRepository {
	return &ProjectRepository{}
}

// Load lee los datos del archivo especificado por PROJECT_FILE_NAME.
func (pr *ProjectRepository) Load() (map[string]interface{}, error) {
	data, err := os.ReadFile(PROJECT_FILE_NAME)
	if err != nil {
		return nil, err
	}

	var project map[string]interface{}
	err = yaml.Unmarshal(data, &project)
	if err != nil {
		return nil, err
	}

	return project, nil
}

// Save escribe un diccionario de datos en el archivo PROJECT_FILE_NAME.
func (pr *ProjectRepository) Save(data map[string]interface{}) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return os.WriteFile(PROJECT_FILE_NAME, yamlData, 0644)
}

package repository

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/workspace/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/workspace/ports"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/workspace/vos"
)

type FileVarsRepository struct{}

func NewFileVarsRepository() ports.VariablesRepository {
	return &FileVarsRepository{}
}

func (r *FileVarsRepository) FindByStep(workspace *aggregates.Workspace, stepFileName vos.FileName) (map[string]string, error) {
	varsFilePath := workspace.VarsFilePath(stepFileName)

	data, err := os.ReadFile(varsFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]string), nil // Devuelve mapa vacío si no existe
		}
		return nil, fmt.Errorf("error reading vars file: %w", err)
	}

	var varsMap map[string]string
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&varsMap); err != nil {
		return nil, fmt.Errorf("error deserializing vars: %w", err)
	}

	return varsMap, nil
}

func (r *FileVarsRepository) SaveByStep(workspace *aggregates.Workspace, stepFileName vos.FileName, vars map[string]string) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(vars); err != nil {
		return fmt.Errorf("error serializing vars to gob: %w", err)
	}

	varsFilePath := workspace.VarsFilePath(stepFileName)

	if err := os.MkdirAll(filepath.Dir(varsFilePath), 0755); err != nil {
		return fmt.Errorf("could not create base directory for vars: %w", err)
	}

	return os.WriteFile(varsFilePath, buffer.Bytes(), 0644)
}

package state

import (
	"encoding/gob"
	"errors"
	"os"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/ports"
)

type GobStateRepository struct{}

func NewGobStateRepository() ports.StateRepository {
	return &GobStateRepository{}
}

func (r *GobStateRepository) Get(filePath string) (*aggregates.StateTable, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Es un caso válido que el archivo no exista la primera vez.
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	var stateTable aggregates.StateTable
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&stateTable); err != nil {
		return nil, err
	}

	return &stateTable, nil
}

func (r *GobStateRepository) Save(filePath string, stateTable *aggregates.StateTable) error {
	// Asegurarse de que el directorio exista
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(stateTable); err != nil {
		return err
	}

	return nil
}

package repository

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/statedetermination/aggregates"
	sdt_ports "github.com/jairoprogramador/fastdeploy-core/internal/domain/statedetermination/ports"
	sdt_vos "github.com/jairoprogramador/fastdeploy-core/internal/domain/statedetermination/vos"
)

type FileStateRepository struct {
	basePath string
}

func NewFileStateRepository(basePath string) sdt_ports.StateRepository {
	return &FileStateRepository{basePath: basePath}
}

func (r *FileStateRepository) Get(workspacePath string, step sdt_vos.Step) (*aggregates.StateTable, error) {
	stateFilePath := r.getStateFilePath(workspacePath, step)

	data, err := os.ReadFile(stateFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			// If file doesn't exist, it's not an error, just means no previous state.
			// The service layer will handle creating a new StateTable.
			return nil, fmt.Errorf("state file not found for step %s: %w", step, err)
		}
		return nil, fmt.Errorf("error reading state file for step %s: %w", step, err)
	}

	var entries []*aggregates.StateEntry
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&entries); err != nil {
		return nil, fmt.Errorf("error deserializing state for step %s: %w", step, err)
	}

	return aggregates.LoadStateTable(step, entries), nil
}

func (r *FileStateRepository) Save(workspacePath string, stateTable *aggregates.StateTable) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(stateTable.Entries()); err != nil {
		return fmt.Errorf("error serializing state for step %s: %w", stateTable.Step(), err)
	}

	stateFilePath := r.getStateFilePath(workspacePath, stateTable.Step())

	if err := os.MkdirAll(filepath.Dir(stateFilePath), 0755); err != nil {
		return fmt.Errorf("could not create base directory for state file: %w", err)
	}

	return os.WriteFile(stateFilePath, buffer.Bytes(), 0644)
}

func (r *FileStateRepository) getStateFilePath(workspacePath string, step sdt_vos.Step) string {
	// e.g., {basePath}/{workspaceName}/.fastdeploy/{step}.state
	fileName := fmt.Sprintf("%s.state", step.String())
	// We use the workspacePath string directly now, assuming it's the root path for the project.
	return filepath.Join(workspacePath, ".fastdeploy", fileName)
}

package repository

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	staAgg "github.com/jairoprogramador/fastdeploy/internal/domain/state/aggregates"
	staPor "github.com/jairoprogramador/fastdeploy/internal/domain/state/ports"
	staVos "github.com/jairoprogramador/fastdeploy/internal/domain/state/vos"

	iStaDto "github.com/jairoprogramador/fastdeploy/internal/infrastructure/state/dto"
	iStaMap "github.com/jairoprogramador/fastdeploy/internal/infrastructure/state/mapper"
)

type StateRepository struct {
	pathStateEnvironment string
	pathStateProject     string
}

func NewStateRepository(
	pathStateRootFastDeploy string,
	projectName string,
	repositoryName string,
	environment string) (staPor.ExecutionStateRepository, error) {

	if pathStateRootFastDeploy == "" {
		return nil, fmt.Errorf("path state root is required")
	}
	if projectName == "" {
		return nil, fmt.Errorf("project name is required")
	}
	if repositoryName == "" {
		return nil, fmt.Errorf("repository name is required")
	}
	if environment == "" {
		return nil, fmt.Errorf("environment is required")
	}

	pathStateEnvironment := filepath.Join(pathStateRootFastDeploy, projectName, repositoryName, environment)
	pathStateProject := filepath.Join(pathStateRootFastDeploy, projectName, repositoryName)
	return &StateRepository{
		pathStateEnvironment: pathStateEnvironment,
		pathStateProject:     pathStateProject,
	}, nil
}

func (r *StateRepository) FindByStepName(stepName string) (*staAgg.ExecutionState, error) {

	fingerprintCode, err := r.findCode()
	if err != nil {
		return nil, err
	}

	fingerprintRecipe, err := r.findRecipe(stepName)
	if err != nil {
		return nil, err
	}

	fingerprintVars, err := r.findVars(stepName)
	if err != nil {
		return nil, err
	}

	dto := iStaDto.StateFingerprintDTO{
		StepName: stepName,
		Fingerprints: map[int]string{
			int(staVos.ScopeCode):   fingerprintCode,
			int(staVos.ScopeRecipe): fingerprintRecipe,
			int(staVos.ScopeVars):   fingerprintVars,
		},
	}

	return iStaMap.ToDomain(dto), nil
}

func (r *StateRepository) Save(state *staAgg.ExecutionState) error {
	stateDTO := iStaMap.ToDTO(state)

	err := r.saveCode(stateDTO.Fingerprints)
	if err != nil {
		return err
	}

	err = r.saveRecipe(stateDTO.StepName, stateDTO.Fingerprints)
	if err != nil {
		return err
	}

	err = r.saveVars(stateDTO.StepName, stateDTO.Fingerprints)
	if err != nil {
		return err
	}

	return nil
}

func (r *StateRepository) saveCode(fingerprints map[int]string) error {
	codeFilePath := r.getPathFileStateCode()
	return r.saveState(codeFilePath, int(staVos.ScopeCode), fingerprints)
}

func (r *StateRepository) saveRecipe(stepName string, fingerprints map[int]string) error {
	recipeFilePath := r.getPathFileStateRecipe(stepName)
	return r.saveState(recipeFilePath, int(staVos.ScopeRecipe), fingerprints)
}

func (r *StateRepository) saveVars(stepName string, fingerprints map[int]string) error {
	varsFilePath := r.getPathFileStateVars(stepName)
	return r.saveState(varsFilePath, int(staVos.ScopeVars), fingerprints)
}

func (r *StateRepository) findCode() (string, error) {
	codeFilePath := r.getPathFileStateCode()
	return r.findState(codeFilePath)
}

func (r *StateRepository) findRecipe(stepName string) (string, error) {
	recipeFilePath := r.getPathFileStateRecipe(stepName)
	return r.findState(recipeFilePath)
}

func (r *StateRepository) findVars(stepName string) (string, error) {
	varsFilePath := r.getPathFileStateVars(stepName)
	return r.findState(varsFilePath)
}

func (r *StateRepository) saveState(filePath string, scope int, fingerprints map[int]string) error {
	fingerprint, ok := fingerprints[scope]

	if !ok {
		return nil
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(fingerprint); err != nil {
		return fmt.Errorf("error al serializar el estado de %s a formato gob: %w", scope, err)
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio base para el estado de %s: %w", scope, err)
	}

	return os.WriteFile(filePath, buffer.Bytes(), 0644)
}

func (r *StateRepository) findState(filePath string) (string, error) {
	dataFile, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("error al leer el archivo de variables: %w", err)
	}

	var fingerprint string
	buffer := bytes.NewBuffer(dataFile)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&fingerprint); err != nil {
		return "", fmt.Errorf("error al deserializar las variables: %w", err)
	}

	return fingerprint, nil
}

func (r *StateRepository) getPathFileStateVars(stepName string) string {
	return filepath.Join(r.pathStateEnvironment, fmt.Sprintf("%s.gob", stepName))
}

func (r *StateRepository) getPathFileStateCode() string {
	return filepath.Join(r.pathStateProject, "code.gob")
}

func (r *StateRepository) getPathFileStateRecipe(stepName string) string {
	return filepath.Join(r.pathStateProject, fmt.Sprintf("%s.gob", stepName))
}

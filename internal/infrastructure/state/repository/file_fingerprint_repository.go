package repository

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	appDto "github.com/jairoprogramador/fastdeploy-core/internal/application/dto"

	shared "github.com/jairoprogramador/fastdeploy-core/internal/domain/shared"
	staAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/aggregates"
	staPor "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/ports"
	staVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"

	iStaDto "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/state/dto"
	iStaMap "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/state/mapper"
)

type FileFingerprintRepository struct {
	pathStateRoot string
}

func NewFileFingerprintRepository(pathStateRoot string) staPor.FingerprintRepository {
	return &FileFingerprintRepository {
		pathStateRoot: pathStateRoot,
	}
}

func (r *FileFingerprintRepository) FindCode(namesRequest appDto.NamesParams) (*staAgg.FingerprintState, error) {
	fingerprintCode, err := r.findCode(namesRequest)
	if err != nil {
		return nil, err
	}

	dto := iStaDto.StateFingerprintDTO{
		StepName:     shared.ScopeCode,
		Fingerprints: map[int]string{int(staVos.ScopeCode): fingerprintCode},
	}

	return iStaMap.ToDomain(dto), nil
}

func (r *FileFingerprintRepository) FindStep(namesRequest appDto.NamesParams, orderRequest appDto.RunParams) (*staAgg.FingerprintState, error) {
	fingerprintRecipe, err := r.findRecipe(namesRequest, orderRequest)
	if err != nil {
		return nil, err
	}

	fingerprintVars, err := r.findVars(namesRequest, orderRequest)
	if err != nil {
		return nil, err
	}

	dto := iStaDto.StateFingerprintDTO{
		StepName: orderRequest.StepName(),
		Fingerprints: map[int]string{
			int(staVos.ScopeRecipe): fingerprintRecipe,
			int(staVos.ScopeVars):   fingerprintVars,
		},
	}

	return iStaMap.ToDomain(dto), nil
}

func (r *FileFingerprintRepository) SaveStep(namesRequest appDto.NamesParams, orderRequest appDto.RunParams, fingerprints *staAgg.FingerprintState) error {
	fingerprintDTO := iStaMap.ToDTO(fingerprints)

	err := r.saveRecipe(namesRequest, orderRequest, fingerprintDTO.Fingerprints)
	if err != nil {
		return err
	}

	err = r.saveVars(namesRequest, orderRequest, fingerprintDTO.Fingerprints)
	if err != nil {
		return err
	}

	return nil
}

func (r *FileFingerprintRepository) SaveCode(namesRequest appDto.NamesParams, fingerprints *staAgg.FingerprintState) error {
	fingerprintDTO := iStaMap.ToDTO(fingerprints)
	codeFilePath := r.getPathFileStateCode(namesRequest)
	return r.saveState(codeFilePath, int(staVos.ScopeCode), fingerprintDTO.Fingerprints)
}

func (r *FileFingerprintRepository) saveRecipe(namesRequest appDto.NamesParams, orderRequest appDto.RunParams, fingerprints map[int]string) error {
	recipeFilePath := r.getPathFileStateRecipe(namesRequest, orderRequest)
	return r.saveState(recipeFilePath, int(staVos.ScopeRecipe), fingerprints)
}

func (r *FileFingerprintRepository) saveVars(namesRequest appDto.NamesParams, orderRequest appDto.RunParams, fingerprints map[int]string) error {
	varsFilePath := r.getPathFileStateVars(namesRequest, orderRequest)
	return r.saveState(varsFilePath, int(staVos.ScopeVars), fingerprints)
}

func (r *FileFingerprintRepository) findCode(namesRequest appDto.NamesParams) (string, error) {
	codeFilePath := r.getPathFileStateCode(namesRequest)
	return r.findState(codeFilePath)
}

func (r *FileFingerprintRepository) findRecipe(namesRequest appDto.NamesParams, orderRequest appDto.RunParams) (string, error) {
	recipeFilePath := r.getPathFileStateRecipe(namesRequest, orderRequest)
	return r.findState(recipeFilePath)
}

func (r *FileFingerprintRepository) findVars(namesRequest appDto.NamesParams, orderRequest appDto.RunParams) (string, error) {
	varsFilePath := r.getPathFileStateVars(namesRequest, orderRequest)
	return r.findState(varsFilePath)
}

func (r *FileFingerprintRepository) saveState(filePath string, scope int, fingerprints map[int]string) error {
	fingerprint, ok := fingerprints[scope]

	if !ok {
		return nil
	}

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(fingerprint); err != nil {
		return fmt.Errorf("error al serializar el estado de %d a formato gob: %w", scope, err)
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("no se pudo crear el directorio base para el estado de %d: %w", scope, err)
	}

	return os.WriteFile(filePath, buffer.Bytes(), 0644)
}

func (r *FileFingerprintRepository) findState(filePath string) (string, error) {
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

func (r *FileFingerprintRepository) getPathFileStateVars(namesRequest appDto.NamesParams, orderRequest appDto.RunParams) string {
	pathStateEnvironment := filepath.Join(r.pathStateRoot, namesRequest.ProjectName(), namesRequest.RepositoryName(), orderRequest.Environment())
	return filepath.Join(pathStateEnvironment, fmt.Sprintf("%s.gob", orderRequest.StepName()))
}

func (r *FileFingerprintRepository) getPathFileStateCode(namesRequest appDto.NamesParams) string {
	pathStateProject := filepath.Join(r.pathStateRoot, namesRequest.ProjectName(), namesRequest.RepositoryName())
	return filepath.Join(pathStateProject, "code.gob")
}

func (r *FileFingerprintRepository) getPathFileStateRecipe(namesRequest appDto.NamesParams, orderRequest appDto.RunParams) string {
	pathStateProject := filepath.Join(r.pathStateRoot, namesRequest.ProjectName(), namesRequest.RepositoryName())
	return filepath.Join(pathStateProject, fmt.Sprintf("%s.gob", orderRequest.StepName()))
}

package adapter

import (
	"deploy/internal/domain/model"
	"fmt"
	"gopkg.in/yaml.v3"
)

type YamlRepository interface {
	Load(pathFile string, out any) model.InfrastructureResponse
	Save(pathFile string, data any) model.InfrastructureResponse
}

type YamlRepositoryImpl struct {
	fileRepository FileRepository
}

func NewYamlRepositoryImpl(fileRepo FileRepository) YamlRepository {
	return &YamlRepositoryImpl{
		fileRepository: fileRepo,
	}
}

func (st *YamlRepositoryImpl) Load(pathFile string, out any) model.InfrastructureResponse {
	file, err := st.fileRepository.OpenFile(pathFile)
	if err != nil {
		return model.NewErrorResponseWithDetails(err, fmt.Sprintf("Failed to open file: %s", pathFile))
	}

	if file == nil {
		return model.NewErrorResponse(fmt.Errorf("file is nil: %s", pathFile))
	}

	defer file.Close()
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(out)
	if err != nil {
		return model.NewErrorResponseWithDetails(err, fmt.Sprintf("Failed to decode YAML from file: %s", pathFile))
	}
	return model.NewResponseWithDetails(out, fmt.Sprintf("Successfully loaded YAML from file: %s", pathFile))
}

func (st *YamlRepositoryImpl) Save(pathFile string, data any) model.InfrastructureResponse {
	file, err := st.fileRepository.CreateFile(pathFile)
	if err != nil {
		return model.NewErrorResponseWithDetails(err, fmt.Sprintf("Failed to create file: %s", pathFile))
	}

	if file == nil {
		return model.NewErrorResponse(fmt.Errorf("file is nil: %s", pathFile))
	}

	defer file.Close()
	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	err = encoder.Encode(data)
	if err != nil {
		return model.NewErrorResponseWithDetails(err, fmt.Sprintf("Failed to encode data to YAML in file: %s", pathFile))
	}

	return model.NewResponseWithDetails(data, fmt.Sprintf("Successfully saved YAML to file: %s", pathFile))
}

package adapter

import (
	"deploy/internal/domain/model"
	"fmt"
	"gopkg.in/yaml.v3"
)

// Error message constants to avoid duplication and improve consistency
const (
	errFileIsNil          = "file is nil: %s"
	errFailedToOpen       = "Failed to open file: %s"
	errFailedToCreate     = "Failed to create file: %s"
	errFailedToDecode     = "Failed to decode YAML from file: %s"
	errFailedToEncode     = "Failed to encode data to YAML in file: %s"
	msgSuccessfullyLoaded = "Successfully loaded YAML from file: %s"
	msgSuccessfullySaved  = "Successfully saved YAML to file: %s"
)

// YamlController defines the interface for YAML file operations
type YamlController interface {
	Load(filePath string, out any) model.InfrastructureResponse
	Save(filePath string, data any) model.InfrastructureResponse
}

// goPkgYamlController implements YamlController using gopkg.in/yaml.v3
type goPkgYamlController struct {
	fileController FileController
}

// NewGoPkgYamlController creates a new YAML controller with the given file controller
func NewGoPkgYamlController(fileController FileController) YamlController {
	return &goPkgYamlController{
		fileController: fileController,
	}
}

// Load reads and decodes YAML from a file into the provided output variable
func (c *goPkgYamlController) Load(filePath string, out any) model.InfrastructureResponse {
	file, err := c.fileController.OpenFile(filePath)
	if err != nil {
		return model.NewErrorResponseWithDetails(err, fmt.Sprintf(errFailedToOpen, filePath))
	}

	if file == nil {
		return model.NewErrorResponse(fmt.Errorf(errFileIsNil, filePath))
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(out); err != nil {
		return model.NewErrorResponseWithDetails(err, fmt.Sprintf(errFailedToDecode, filePath))
	}

	return model.NewResponseWithDetails(out, fmt.Sprintf(msgSuccessfullyLoaded, filePath))
}

// Save encodes and writes data as YAML to a file
func (c *goPkgYamlController) Save(filePath string, data any) model.InfrastructureResponse {
	file, err := c.fileController.CreateFile(filePath)
	if err != nil {
		return model.NewErrorResponseWithDetails(err, fmt.Sprintf(errFailedToCreate, filePath))
	}

	if file == nil {
		return model.NewErrorResponse(fmt.Errorf(errFileIsNil, filePath))
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	if err := encoder.Encode(data); err != nil {
		return model.NewErrorResponseWithDetails(err, fmt.Sprintf(errFailedToEncode, filePath))
	}

	return model.NewResponseWithDetails(data, fmt.Sprintf(msgSuccessfullySaved, filePath))
}

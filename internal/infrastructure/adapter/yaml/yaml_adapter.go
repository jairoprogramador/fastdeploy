package yaml

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/adapter/file"
	"github.com/jairoprogramador/fastdeploy/pkg/common/logger"
	"gopkg.in/yaml.v3"
)

const (
	errFileIsNil      = "file is nil: %s"
	errFailedToDecode = "failed to decode YAML from file: %s, the error is %v"
	errFailedToEncode = "failed to encode data to YAML in file: %s, the error is %v"
)

type YamlPort interface {
	Load(filePath string, out any) error
	Save(filePath string, data any) error
}

type yamlAdapter struct {
	filePort   file.FilePort
	fileLogger *logger.FileLogger
}

func NewYamlAdapter(filePort file.FilePort, fileLogger *logger.FileLogger) YamlPort {
	return &yamlAdapter{
		filePort:   filePort,
		fileLogger: fileLogger,
	}
}

func (c *yamlAdapter) Load(filePath string, out any) error {
	file, err := c.filePort.OpenFile(filePath)
	if err != nil {
		return err
	}

	if file == nil {
		return c.logError(fmt.Errorf(errFileIsNil, filePath))
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(out); err != nil {
		return c.logError(fmt.Errorf(errFailedToDecode, filePath, err))
	}

	return nil
}

func (c *yamlAdapter) Save(filePath string, data any) error {
	file, err := c.filePort.CreateFile(filePath)
	if err != nil {
		return err
	}

	if file == nil {
		return c.logError(fmt.Errorf(errFileIsNil, filePath))
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	defer encoder.Close()

	if err := encoder.Encode(data); err != nil {
		return c.logError(fmt.Errorf(errFailedToEncode, filePath, err))
	}

	return nil
}

func (c *yamlAdapter) logError(err error) error {
	if err != nil {
		c.fileLogger.Error(err)
	}
	return err
}

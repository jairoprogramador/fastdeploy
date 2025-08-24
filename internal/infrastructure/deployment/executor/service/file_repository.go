package service

import (
	"os"
	"fmt"

	"gopkg.in/yaml.v3"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/dto"
)

func Load(filePath string) (dto.ListCmdDto, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return dto.ListCmdDto{}, fmt.Errorf("%w: config file does not exist in " + filePath, err)
		}
		return dto.ListCmdDto{}, err
	}

	var result dto.ListCmdDto
	if err := yaml.Unmarshal(data, &result); err != nil {
		return dto.ListCmdDto{}, err
	}

	return result, nil
}

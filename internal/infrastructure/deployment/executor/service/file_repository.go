package service

import (
	"os"

	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/dto"
	"gopkg.in/yaml.v3"
)

func LoadCmdList(filePath string) (dto.CmdListDTO, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return dto.CmdListDTO{}, nil
		}
		return dto.CmdListDTO{}, err
	}

	var result dto.CmdListDTO
	if err := yaml.Unmarshal(data, &result); err != nil {
		return dto.CmdListDTO{}, err
	}

	return result, nil
}

func LoadVariableList(filePath string) (dto.VariableListDTO, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return dto.VariableListDTO{}, nil
		}
		return dto.VariableListDTO{}, err
	}

	var result dto.VariableListDTO
	if err := yaml.Unmarshal(data, &result); err != nil {
		return dto.VariableListDTO{}, err
	}

	return result, nil
}

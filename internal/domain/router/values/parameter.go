package values

import (
	"errors"
	"strings"
)

type Parameter struct {
	homeDir         string
	projectName     string
	repositoryName  string
	environment     string
	stackName       string
}

func NewParameter(
	homeDir string,
	projectName string,
	repositoryName string,
	environment string,
	stackName string,
) (Parameter, error) {

	homeDir, err := validateIsEmpty(homeDir)
	if err != nil {
		return Parameter{}, err
	}

	projectName, err = validateIsEmpty(projectName)
	if err != nil {
		return Parameter{}, err
	}

	repositoryName, err = validateIsEmpty(repositoryName)
	if err != nil {
		return Parameter{}, err
	}

	environment, err = validateIsEmpty(environment)
	if err != nil {
		return Parameter{}, err
	}

	stackName = strings.TrimSpace(stackName)

	return Parameter{
		homeDir:         homeDir,
		projectName:     projectName,
		repositoryName:  repositoryName,
		environment:     environment,
		stackName:       stackName,
	}, nil
}

func validateIsEmpty(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errors.New("value cannot be empty")
	}
	return value, nil
}

func (p Parameter) GetHomeDir() string {
	return p.homeDir
}

func (p Parameter) GetProjectName() string {
	return p.projectName
}

func (p Parameter) GetRepositoryName() string {
	return p.repositoryName
}

func (p Parameter) GetEnvironment() string {
	return p.environment
}

func (p Parameter) GetStackName() string {
	return p.stackName
}

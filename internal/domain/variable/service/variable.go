package service

import (
	"fmt"
	"log"
	"regexp"

	context "github.com/jairoprogramador/fastdeploy/internal/domain/context/values"
	"github.com/jairoprogramador/fastdeploy/internal/domain/router/service"
	"github.com/jairoprogramador/fastdeploy/internal/domain/variable/port"
)

const COMPUTED_DIR_NAME = "variables"
const COMPUTED_FILE_NAME = "computed.yaml"

var varRegexFormat = regexp.MustCompile(`\$\{var\.([^}]+)\}`)

type VariableService interface {
	Process(value string, context *context.ContextValue) string
	AddVariablesComputed(stepName string, context *context.ContextValue) (*context.ContextValue, error)
	AddVariablesStep(stepName string, context *context.ContextValue) (*context.ContextValue, error)
}

type VariableServiceImpl struct {
	variablePort   port.VariablePort
	pathPepository service.RepositoryRouterService
}

func NewVariableService(
	variablePort port.VariablePort,
	pathPepository service.RepositoryRouterService) VariableService {

	return &VariableServiceImpl{
		variablePort:   variablePort,
		pathPepository: pathPepository,
	}
}

func (v *VariableServiceImpl) AddVariablesComputed(
	stepName string,
	context *context.ContextValue) (*context.ContextValue, error) {

	pathStep := v.pathPepository.GetPathStep(stepName)
	pathFile := v.pathPepository.BuildPath(pathStep, COMPUTED_DIR_NAME, COMPUTED_FILE_NAME)

	return v.addVariables(pathFile, context)
}

func (v *VariableServiceImpl) AddVariablesStep(
	stepName string,
	context *context.ContextValue) (*context.ContextValue, error) {

	nameFile := fmt.Sprintf("%s.yaml", stepName)

	pathStep := v.pathPepository.GetPathStep(stepName)
	pathFile := v.pathPepository.BuildPath(pathStep, COMPUTED_DIR_NAME, nameFile)

	return v.addVariables(pathFile, context)
}

func (v *VariableServiceImpl) addVariables(
	pathFile string,
	context *context.ContextValue) (*context.ContextValue, error) {

	variables, err := v.variablePort.Load(pathFile)
	if err != nil {
		return context, err
	}

	for _, variable := range variables {
		value := v.Process(variable.GetValue(), context)
		context.Set(variable.GetName(), value)
	}

	return context, nil
}

func (v *VariableServiceImpl) Process(
	value string,
	context *context.ContextValue) string {

	return varRegexFormat.ReplaceAllStringFunc(
		value,
		func(match string) string {
			subMatch := varRegexFormat.FindStringSubmatch(match)
			if len(subMatch) >= 1 {
				value, err := context.Get(subMatch[1])
				if err != nil {
					log.Printf("error al obtener el valor de la variable %s: %v", subMatch[1], err)
					return match
				}
				if value != "" {
					return value
				}
			}
			log.Printf("no se encontró coincidencia para la variable %s", match)
			return match
		},
	)
}

package project

import (
	"fmt"
	"strings"
)

type ProjectValidatorImpl struct{}

func NewProjectValidator() ProjectValidator {
	return &ProjectValidatorImpl{}
}

func (dpv *ProjectValidatorImpl) Validate(project ProjectEntity) error {
	var errors []string

	if err := dpv.validateRequiredFields(project); err != nil {
		errors = append(errors, err.Error())
	}

	if err := dpv.validateFormats(project); err != nil {
		errors = append(errors, err.Error())
	}

	if err := dpv.validateBusinessRules(project); err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return fmt.Errorf("errores de validación: %s", strings.Join(errors, "; "))
	}

	return nil
}

func (dpv *ProjectValidatorImpl) validateRequiredFields(project ProjectEntity) error {
	if project.ProjectName == "" {
		return fmt.Errorf("el nombre del proyecto es requerido")
	}
	if project.Organization == "" {
		return fmt.Errorf("la organización es requerida")
	}
	if project.Technology == "" {
		return fmt.Errorf("la tecnología es requerida")
	}
	return nil
}

func (dpv *ProjectValidatorImpl) validateFormats(project ProjectEntity) error {
	if project.ProjectID != "" && len(project.ProjectID) != 40 {
		return fmt.Errorf("el ProjectID debe tener 40 caracteres (SHA1)")
	}
	return nil
}

func (dpv *ProjectValidatorImpl) validateBusinessRules(project ProjectEntity) error {
	supportedTechnologies := []string{"springboot", "node", "python", "go"}
	isSupported := false
	for _, tech := range supportedTechnologies {
		if project.Technology == tech {
			isSupported = true
			break
		}
	}
	if !isSupported {
		return fmt.Errorf("tecnología no soportada: %s", project.Technology)
	}
	return nil
}

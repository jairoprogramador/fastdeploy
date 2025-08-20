package services

import (
	"fmt"
	"strings"

	"github.com/jairoprogramador/fastdeploy/internal/domain/project/entities"
)

type Validator interface {
	Validate(project entities.Project) error
}

type ValidatorProject struct{}

func NewValidatorProject() Validator {
	return &ValidatorProject{}
}

func (dpv *ValidatorProject) Validate(project entities.Project) error {
	var errors []string

	if err := dpv.validateRequiredFields(project); err != nil {
		errors = append(errors, err.Error())
	}

	if err := dpv.validateFormats(project); err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return fmt.Errorf("errores de validación: %s", strings.Join(errors, "; "))
	}

	return nil
}


func (pv *ValidatorProject) validateRequiredFields(project entities.Project) error {
	var validationErrors []string

	if project.GetID().IsEmpty() {
		validationErrors = append(validationErrors, "project id is empty")
	}

	if project.GetName().IsEmpty() {
		validationErrors = append(validationErrors, "project name is empty")
	}

	if project.GetOrganization().IsEmpty() {
		validationErrors = append(validationErrors, "organization is empty")
	}

	if project.GetRepository().GetURL().IsEmpty() {
		validationErrors = append(validationErrors, "repository url is empty")
	}

	if project.GetTechnology().GetName().IsEmpty() {
		validationErrors = append(validationErrors, "technology name is empty")
	}

	if project.GetTechnology().GetVersion().IsEmpty() {
		validationErrors = append(validationErrors, "technology version is empty")
	}

	if project.GetDeployment().GetVersion().IsEmpty() {
		validationErrors = append(validationErrors, "deployment version is empty")
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("project validation failed: %v", validationErrors)
	}

	return nil
}

func (dpv *ValidatorProject) validateFormats(project entities.Project) error {
	if !project.GetID().IsValid() {
		return fmt.Errorf("el ProjectID debe tener 16 caracteres como mínimo")
	}

	if !project.GetName().IsValid() {
		return fmt.Errorf("el nombre del proyecto no es válido")
	}

	if !project.GetOrganization().IsValid() {
		return fmt.Errorf("el nombre de la organización no es válido")
	}

	if !project.GetRepository().GetURL().IsValid() {
		return fmt.Errorf("el repositorio no es válido")
	}

	if !project.GetTechnology().GetName().IsValid() {
		return fmt.Errorf("el nombre de la tecnología no es válido")
	}

	if !project.GetTechnology().GetVersion().IsValid() {
		return fmt.Errorf("la versión de la tecnología no es válida")
	}

	if !project.GetDeployment().GetVersion().IsValid() {
		return fmt.Errorf("la versión del despliegue no es válida")
	}
	
	return nil
}


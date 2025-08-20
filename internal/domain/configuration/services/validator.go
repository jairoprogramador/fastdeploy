package services


import (
	"fmt"
	"strings"

	"github.com/jairoprogramador/fastdeploy/internal/domain/configuration/entities"
)

type Validator interface {
	Validate(config entities.Configuration) error
}

type ValidatorConfiguration struct{}

func NewValidatorConfiguration() Validator {
	return &ValidatorConfiguration{}
}

func (dpv *ValidatorConfiguration) Validate(config entities.Configuration) error {
	var errors []string

	if err := dpv.validateRequiredFields(config); err != nil {
		errors = append(errors, err.Error())
	}

	if err := dpv.validateFormats(config); err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return fmt.Errorf("errores de validación: %s", strings.Join(errors, "; "))
	}

	return nil
}


func (pv *ValidatorConfiguration) validateRequiredFields(config entities.Configuration) error {
	var validationErrors []string

	if config.GetNameOrganization().IsEmpty() {
		validationErrors = append(validationErrors, "name organization is empty")
	}

	if config.GetRepository().GetURL().IsEmpty() {
		validationErrors = append(validationErrors, "repository url is empty")
	}

	if config.GetTechnology().GetName().IsEmpty() {
		validationErrors = append(validationErrors, "technology name is empty")
	}

	if config.GetTechnology().GetVersion().IsEmpty() {
		validationErrors = append(validationErrors, "technology version is empty")
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf("project validation failed: %v", validationErrors)
	}

	return nil
}

func (dpv *ValidatorConfiguration) validateFormats(config entities.Configuration) error {
	if !config.GetNameOrganization().IsValid() {
		return fmt.Errorf("el nombre de la organización no es válido")
	}

	if !config.GetRepository().GetURL().IsValid() {
		return fmt.Errorf("el repositorio no es válido")
	}

	if !config.GetTechnology().GetName().IsValid() {
		return fmt.Errorf("el nombre de la tecnología no es válido")
	}

	if !config.GetTechnology().GetVersion().IsValid() {
		return fmt.Errorf("la versión de la tecnología no es válida")
	}
	
	return nil
}

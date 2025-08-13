package config

import (
	"fmt"
	"strings"
)

type ConfigValidatorImpl struct{}

func NewConfigValidatorImpl() ConfigValidator {
	return &ConfigValidatorImpl{}
}

func (dcv *ConfigValidatorImpl) Validate(config ConfigEntity) error {
	var errors []string

	// Validar campos requeridos
	if err := dcv.validateRequiredFields(config); err != nil {
		errors = append(errors, err.Error())
	}

	// Validar formatos
	if err := dcv.validateFormats(config); err != nil {
		errors = append(errors, err.Error())
	}

	// Validar reglas de negocio
	if err := dcv.validateBusinessRules(config); err != nil {
		errors = append(errors, err.Error())
	}

	if len(errors) > 0 {
		return fmt.Errorf("errores de validación: %s", strings.Join(errors, "; "))
	}

	return nil
}

func (dcv *ConfigValidatorImpl) validateRequiredFields(config ConfigEntity) error {
	// Implementar validaciones específicas según tu ConfigEntity
	// Por ejemplo:
	// if config.Name == "" {
	//     return fmt.Errorf("el nombre es requerido")
	// }
	return nil
}

func (dcv *ConfigValidatorImpl) validateFormats(config ConfigEntity) error {
	// Implementar validaciones de formato
	// Por ejemplo, validar emails, URLs, etc.
	return nil
}

func (dcv *ConfigValidatorImpl) validateBusinessRules(config ConfigEntity) error {
	// Implementar reglas de negocio específicas
	return nil
}

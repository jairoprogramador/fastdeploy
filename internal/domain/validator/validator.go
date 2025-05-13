package validator

import (
	"deploy/internal/domain/condition"
	"deploy/internal/domain/model"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/locales/es"
	translator "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	dictionary_es "github.com/go-playground/validator/v10/translations/es"
)

const (
	//TypeCheck   = "check"
	TypeCommand   = "command"
	TypeContainer = "container"
	//TypeHttp    = "http"
)

type DeploymentValidator struct {
	validate *validator.Validate
	trans    translator.Translator
}

var (
	instance                *DeploymentValidator
	onceDeploymentValidator sync.Once
)

func GetDeploymentValidator() *DeploymentValidator {
	onceDeploymentValidator.Do(func() {
		locatorEs := es.New()
		universalTranslator := translator.New(locatorEs, locatorEs)
		translatorEs, _ := universalTranslator.GetTranslator("es")

		validatorStruct := validator.New()
		_ = dictionary_es.RegisterDefaultTranslations(validatorStruct, translatorEs)

		instance = &DeploymentValidator{
			validate: validatorStruct,
			trans:    translatorEs,
		}
	})
	return instance
}

func (v *DeploymentValidator) Validate(deployment *model.Deployment) error {
	if err := v.validate.Struct(deployment); err != nil {
		return v.translateError(err)
	}
	return v.validateSteps(deployment.Steps)
}

func (v *DeploymentValidator) validateSteps(steps []model.Step) error {
	stepNames := make(map[string]bool)

	allStepNames := make(map[string]bool)
	for _, step := range steps {
		allStepNames[step.Name] = true
	}

	for _, step := range steps {

		if step.Name == "" {
			return fmt.Errorf("nombre de paso invalido: %s", step.Name)
		}

		if exists := stepNames[step.Name]; exists {
			return fmt.Errorf("nombre de paso duplicado: %s", step.Name)
		}
		stepNames[step.Name] = true

		if step.Type == "" {
			step.Type = TypeCommand
		}

		if !v.isValidStepType(step.Type) {
			return fmt.Errorf("tipo inválido: %s en paso %s", step.Type, step.Name)
		}

		// Validar dependencias
		/* if err := v.validateDependencies(step, stepNames); err != nil {
			return err
		} */

		if err := v.validateConditions(step, allStepNames); err != nil {
			return err
		}

		if step.Timeout != "" {
			if _, err := time.ParseDuration(step.Timeout); err != nil {
				return fmt.Errorf("formato de timeout inválido en paso %s: %s", step.Name, step.Timeout)
			}
		}

		if step.Retry != nil {
			if step.Retry.Attempts < 1 {
				return fmt.Errorf("número de intentos inválido en paso %s", step.Name)
			}
			if _, err := time.ParseDuration(step.Retry.Delay); err != nil {
				return fmt.Errorf("formato de delay inválido en paso %s: %s", step.Name, step.Retry.Delay)
			}
		}
	}

	return nil
}

func (v *DeploymentValidator) isValidStepType(stepType string) bool {
	validTypes := map[string]bool{
		//TypeCheck:   true,
		TypeCommand: true,
		//TypeHttp:    true,
	}
	return validTypes[stepType]
}

/*
	 func (v *DeploymentValidator) validateDependencies(step model.Step, existingSteps map[string]bool) error {
		for _, dep := range step.DependsOn {
			if !existingSteps[dep] {
				return fmt.Errorf("dependencia no encontrada en paso %s: %s", step.Name, dep)
			}
		}
		return nil
	}
*/
func (v *DeploymentValidator) validateConditions(step model.Step, stepNames map[string]bool) error {
	if step.If == "" {
		return nil
	}

	parts := strings.SplitN(step.If, ":", 2)
	if len(parts) == 0 {
		return fmt.Errorf("formato de condición inválido en paso %s", step.Name)
	}
	
	validConditions := map[string]bool{
		string(condition.NotEmpty): true,
		string(condition.Empty):    true,
		string(condition.Equals):   true,
		string(condition.Contains): true,
		string(condition.Matches):  true,
	}

	if exists := validConditions[parts[0]]; !exists {
		return fmt.Errorf("tipo de condición inválido en paso %s: %s", step.Name, parts[0])
	}

	if parts[0] == string(condition.Equals) || parts[0] == string(condition.Contains) || parts[0] == string(condition.Matches) {
		if len(parts) != 2 {
			return fmt.Errorf("valor requerido para condición %s en paso %s", parts[0], step.Name)
		}
	}

	if step.Then != "" && !stepNames[step.Then] {
		return fmt.Errorf("paso destino no encontrado en condición de %s: %s", step.Name, step.Then)
	}

	return nil
}

func (v *DeploymentValidator) translateError(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrors {
			return fmt.Errorf("%s", e.Translate(v.trans))
		}
	}
	return err
}

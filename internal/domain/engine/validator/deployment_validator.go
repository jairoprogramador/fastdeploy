package validator

import (
	"deploy/internal/domain/engine/condition"
	"deploy/internal/domain/model"
	"deploy/internal/domain/model/logger"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/locales/es"
	translator "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	dictionary_es "github.com/go-playground/validator/v10/translations/es"
)

const (
	TypeCommand   = "command"
	TypeContainer = "container"
	TypeSetup     = "setup"
)

const (
	ThenFinish = "finish"
)

const (
	prefix = "invalid deployment"
)

type DeploymentValidator struct {
	validate *validator.Validate
	trans    translator.Translator
	logger   *logger.Logger
}

func NewDeploymentValidator(logger *logger.Logger) *DeploymentValidator {
	locatorEs := es.New()
	universalTranslator := translator.New(locatorEs, locatorEs)
	translatorEs, _ := universalTranslator.GetTranslator("es")

	validatorStruct := validator.New()
	_ = dictionary_es.RegisterDefaultTranslations(validatorStruct, translatorEs)

	return &DeploymentValidator{
		validate: validatorStruct,
		trans:    translatorEs,
		logger:   logger,
	}
}

func (v *DeploymentValidator) Validate(deployment *model.DeploymentEntity) error {
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
			message := fmt.Sprintf("%s: step name is empty", prefix)
			return v.logger.NewError(message)
		}

		if exists := stepNames[step.Name]; exists {
			message := fmt.Sprintf("%s: duplicate step name: '%s'", prefix, step.Name)
			return v.logger.NewError(message)
		}
		stepNames[step.Name] = true

		if !v.isValidStepType(step.Type) {
			message := fmt.Sprintf("%s: step type: '%s' in step %s", prefix, step.Type, step.Name)
			return v.logger.NewError(message)
		}

		if err := v.validateConditions(step, allStepNames); err != nil {
			return err
		}

		if step.Timeout != "" {
			if _, err := time.ParseDuration(step.Timeout); err != nil {
				message := fmt.Sprintf("%s: timeout format: '%s' in step %s", prefix, step.Timeout, step.Name)
				return v.logger.NewError(message)
			}
		}

		if step.Retry != nil {
			if step.Retry.Attempts < 1 {
				message := fmt.Sprintf("%s: incorrect number of attempts: %d in step %s", prefix, step.Retry.Attempts, step.Name)
				return v.logger.NewError(message)
			}
			if _, err := time.ParseDuration(step.Retry.Delay); err != nil {
				message := fmt.Sprintf("%s: delay format: '%s' in step %s", prefix, step.Retry.Delay, step.Name)
				return v.logger.NewError(message)
			}
		}
	}
	return nil
}

func (v *DeploymentValidator) isValidStepType(stepType string) bool {
	validTypes := map[string]bool{
		TypeContainer: true,
		TypeCommand:   true,
		TypeSetup:     true,
	}
	return validTypes[stepType]
}

func (v *DeploymentValidator) validateConditions(step model.Step, stepNames map[string]bool) error {
	if step.If == "" {
		return nil
	}

	parts := strings.SplitN(step.If, ":", 2)
	if len(parts) == 0 {
		message := fmt.Sprintf("%s: invalid condition format in step %s", prefix, step.Name)
		return v.logger.NewError(message)
	}

	validConditions := map[string]bool{
		string(condition.NotEmpty): true,
		string(condition.Empty):    true,
		string(condition.Equals):   true,
		string(condition.Contains): true,
		string(condition.Matches):  true,
	}

	if exists := validConditions[parts[0]]; !exists {
		message := fmt.Sprintf("%s: invalid condition type: '%s' in step %s", prefix, parts[0], step.Name)
		return v.logger.NewError(message)
	}

	if parts[0] == string(condition.Equals) || parts[0] == string(condition.Contains) || parts[0] == string(condition.Matches) {
		if len(parts) != 2 {
			message := fmt.Sprintf("%s: required value for condition %s in step %s", prefix, parts[0], step.Name)
			return v.logger.NewError(message)
		}
		if parts[1] == "" {
			message := fmt.Sprintf("%s: required value for condition %s in step %s", prefix, parts[0], step.Name)
			return v.logger.NewError(message)
		}
	}

	if parts[0] == string(condition.Matches) {
		if _, err := regexp.Compile(parts[1]); err != nil {
			message := fmt.Sprintf("%s: invalid regex pattern: '%s' in step %s", prefix, parts[1], step.Name)
			return v.logger.NewError(message)
		}
	}

	if step.Then != "" && !stepNames[step.Then] {
		message := fmt.Sprintf("%s: destination step not found in condition of %s: %s", prefix, step.Name, step.Then)
		return v.logger.NewError(message)
	}

	return nil
}

func (v *DeploymentValidator) translateError(err error) error {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			return fmt.Errorf("%s", e.Translate(v.trans))
		}
	}
	return err
}

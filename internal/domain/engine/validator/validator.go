package validator

import (
	"errors"
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/model"
	"regexp"
	"strings"
	"time"

	"github.com/go-playground/locales/es"
	translator "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	dictionary_es "github.com/go-playground/validator/v10/translations/es"
	"github.com/jairoprogramador/fastdeploy/internal/domain/engine/condition"
)

var validConditionTypes = map[string]bool{
	string(condition.NotEmpty): true,
	string(condition.Empty):    true,
	string(condition.Equals):   true,
	string(condition.Contains): true,
	string(condition.Matches):  true,
}

var validStepTypes = map[string]bool{
	string(model.Container): true,
	string(model.Command):   true,
	string(model.Check):     true,
}

const (
	ThenFinish = "finish"
)

const (
	ErrorPrefix           = "invalid deployment"
	ErrorEmptyStepName    = "%s: step name is empty"
	ErrorDuplicateStep    = "%s: duplicate step name: '%s'"
	ErrorInvalidStepType  = "%s: step type: '%s' in step %s"
	ErrorInvalidCondition = "%s: condition format in step %s"
	ErrorInvalidCondType  = "%s: condition type: '%s' in step %s"
	ErrorRequiredValue    = "%s: required value for condition %s in step %s"
	ErrorInvalidRegex     = "%s: regex pattern: '%s' in step %s"
	ErrorDestStepNotFound = "%s: destination step not found in condition of %s: %s"
	ErrorTimeoutFormat    = "%s: timeout format: '%s' in step %s"
	ErrorRetryAttempts    = "%s: incorrect number of attempts: %d in step %s"
	ErrorRetryDelay       = "%s: delay format: '%s' in step %s"
)

const (
	ConditionSeparator = ":"
	MinRetryAttempts   = 1
)

type Validator struct {
	validate *validator.Validate
	trans    translator.Translator
}

func NewValidator() *Validator {
	locatorEs := es.New()
	universalTranslator := translator.New(locatorEs, locatorEs)
	translatorEs, _ := universalTranslator.GetTranslator("es")

	validatorStruct := validator.New()
	_ = dictionary_es.RegisterDefaultTranslations(validatorStruct, translatorEs)

	return &Validator{
		validate: validatorStruct,
		trans:    translatorEs,
	}
}

func (v *Validator) Validate(deployment *model.DeploymentEntity) error {
	if err := v.validate.Struct(deployment); err != nil {
		return v.translateError(err)
	}
	return v.validateSteps(deployment.Steps)
}

func (v *Validator) validateSteps(steps []model.Step) error {
	stepNames := make(map[string]bool)
	allStepNames := collectStepNames(steps)

	for _, step := range steps {
		if err := v.validateStepName(step.Name, stepNames); err != nil {
			return err
		}

		if err := v.validateStepType(step); err != nil {
			return err
		}

		if err := v.validateConditions(step, allStepNames); err != nil {
			return err
		}

		if err := v.validateTimeout(step); err != nil {
			return err
		}

		if err := v.validateRetry(step); err != nil {
			return err
		}

		stepNames[step.Name] = true
	}
	return nil
}

func collectStepNames(steps []model.Step) map[string]bool {
	names := make(map[string]bool)
	for _, step := range steps {
		names[step.Name] = true
	}
	return names
}

func (v *Validator) validateStepName(name string, existingNames map[string]bool) error {
	if name == "" {
		return fmt.Errorf(ErrorEmptyStepName, ErrorPrefix)
	}

	if existingNames[name] {
		return fmt.Errorf(ErrorDuplicateStep, ErrorPrefix, name)
	}
	return nil
}

func (v *Validator) validateStepType(step model.Step) error {
	if !validStepTypes[step.Type] {
		return fmt.Errorf(ErrorInvalidStepType, ErrorPrefix, step.Type, step.Name)
	}
	return nil
}

func (v *Validator) validateTimeout(step model.Step) error {
	if step.Timeout == "" {
		return nil
	}

	if _, err := time.ParseDuration(step.Timeout); err != nil {
		return fmt.Errorf(ErrorTimeoutFormat, ErrorPrefix, step.Timeout, step.Name)
	}
	return nil
}

func (v *Validator) validateRetry(step model.Step) error {
	if step.Retry == nil {
		return nil
	}

	if step.Retry.Attempts < MinRetryAttempts {
		return fmt.Errorf(ErrorRetryAttempts, ErrorPrefix, step.Retry.Attempts, step.Name)
	}

	if _, err := time.ParseDuration(step.Retry.Delay); err != nil {
		return fmt.Errorf(ErrorRetryDelay, ErrorPrefix, step.Retry.Delay, step.Name)
	}

	return nil
}

func (v *Validator) validateConditions(step model.Step, allStepNames map[string]bool) error {
	if step.If == "" {
		return nil
	}

	parts := strings.SplitN(step.If, ConditionSeparator, 2)
	if len(parts) == 0 {
		return fmt.Errorf(ErrorInvalidCondition, ErrorPrefix, step.Name)
	}

	conditionType := parts[0]
	if err := v.validateConditionType(conditionType, step.Name); err != nil {
		return err
	}

	if err := v.validateConditionValue(conditionType, parts, step.Name); err != nil {
		return err
	}

	if err := v.validateRegexPattern(conditionType, parts, step.Name); err != nil {
		return err
	}

	return v.validateThenStep(step, allStepNames)
}

func (v *Validator) validateConditionType(condType string, stepName string) error {
	if !validConditionTypes[condType] {
		return fmt.Errorf(ErrorInvalidCondType, ErrorPrefix, condType, stepName)
	}

	return nil
}

func (v *Validator) validateConditionValue(condType string, parts []string, stepName string) error {
	needsValue := condType == string(condition.Equals) ||
		condType == string(condition.Contains) ||
		condType == string(condition.Matches)

	if !needsValue {
		return nil
	}

	if len(parts) != 2 || parts[1] == "" {
		return fmt.Errorf(ErrorRequiredValue, ErrorPrefix, condType, stepName)
	}

	return nil
}

func (v *Validator) validateRegexPattern(condType string, parts []string, stepName string) error {
	if condType != string(condition.Matches) || len(parts) < 2 {
		return nil
	}

	if _, err := regexp.Compile(parts[1]); err != nil {
		return fmt.Errorf(ErrorInvalidRegex, ErrorPrefix, parts[1], stepName)
	}

	return nil
}

func (v *Validator) validateThenStep(step model.Step, allStepNames map[string]bool) error {
	if step.Then == "" || allStepNames[step.Then] {
		return nil
	}

	return fmt.Errorf(ErrorDestStepNotFound, ErrorPrefix, step.Name, step.Then)
}

func (v *Validator) translateError(err error) error {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, e := range validationErrors {
			return fmt.Errorf("%s", e.Translate(v.trans))
		}
	}
	return err
}

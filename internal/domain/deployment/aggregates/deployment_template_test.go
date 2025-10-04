package aggregates

import (
	"testing"

	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entities"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
)

const ENV_STAGING = "staging"
const STEP_TEST = "test"
// --- Helper para crear datos de prueba ---
func createValidTemplateSource(t *testing.T) vos.TemplateSource {
	t.Helper()
	source, err := vos.NewTemplateSource("http://test.com/repo.git", "main")
	if err != nil {
		t.Fatalf("fallo al crear helper TemplateSource: %v", err)
	}
	return source
}

func createValidEnvironments(t *testing.T) []vos.Environment {
	t.Helper()
	env1, err1 := vos.NewEnvironment(ENV_STAGING, "Staging Env", "stag")
	env2, err2 := vos.NewEnvironment("production", "Production Env", "prod")
	if err1 != nil || err2 != nil {
		t.Fatalf("fallo al crear helpers Environment: %v, %v", err1, err2)
	}
	return []vos.Environment{env1, env2}
}

func createValidSteps(t *testing.T) []entities.StepDefinition {
	t.Helper()
	verifications := []vos.VerificationType{vos.VerificationTypeCode}
	verifications2 := []vos.VerificationType{vos.VerificationTypeEnv}
	cmd, err := vos.NewCommandDefinition("test-cmd", "echo")
	if err != nil {
		t.Fatalf("fallo al crear helper CommandDefinition: %v", err)
	}
	step1, err1 := entities.NewStepDefinition(STEP_TEST, verifications, []vos.CommandDefinition{cmd})
	step2, err2 := entities.NewStepDefinition("deploy", verifications2, []vos.CommandDefinition{cmd})
	if err1 != nil || err2 != nil {
		t.Fatalf("fallo al crear helpers StepDefinition: %v, %v", err1, err2)
	}
	return []entities.StepDefinition{step1, step2}
}

// --- Tests ---
func TestNewDeploymentTemplate(t *testing.T) {
	validSource := createValidTemplateSource(t)
	validEnvs := createValidEnvironments(t)
	validSteps := createValidSteps(t)

	verifications := []vos.VerificationType{vos.VerificationTypeCode}

	// Crear un environment duplicado para el caso de prueba de fallo
	envDupe, _ := vos.NewEnvironment(ENV_STAGING, "Duplicated Staging", "stag-dupe")
	dupeEnvs := append(validEnvs, envDupe)

	// Crear un paso duplicado
	cmd, _ := vos.NewCommandDefinition("cmd", "c")
	stepDupe, _ := entities.NewStepDefinition(STEP_TEST, verifications, []vos.CommandDefinition{cmd})
	dupeSteps := append(validSteps, stepDupe)

	testCases := []struct {
		testName     string
		source       vos.TemplateSource
		environments []vos.Environment
		steps        []entities.StepDefinition
		expectError  bool
	}{
		{
			testName:     "Creacion valida",
			source:       validSource,
			environments: validEnvs,
			steps:        validSteps,
			expectError:  false,
		},
		{
			testName:     "Fallo por environments vacios",
			source:       validSource,
			environments: []vos.Environment{},
			steps:        validSteps,
			expectError:  true,
		},
		{
			testName:     "Fallo por steps vacios",
			source:       validSource,
			environments: validEnvs,
			steps:        []entities.StepDefinition{},
			expectError:  true,
		},
		{
			testName:     "Fallo por environments duplicados",
			source:       validSource,
			environments: dupeEnvs,
			steps:        validSteps,
			expectError:  true,
		},
		{
			testName:     "Fallo por steps duplicados",
			source:       validSource,
			environments: validEnvs,
			steps:        dupeSteps,
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			_, err := NewDeploymentTemplate(tc.source, tc.environments, tc.steps)

			if tc.expectError && err == nil {
				t.Errorf("Se esperaba un error, pero no se obtuvo ninguno")
			}
			if !tc.expectError && err != nil {
				t.Errorf("No se esperaba un error, pero se obtuvo: %v", err)
			}
		})
	}
}

func TestDeploymentTemplate_DefensiveCopying(t *testing.T) {
	t.Run("Environments debe devolver una copia", func(t *testing.T) {
		originalSteps := createValidSteps(t)
		originalEnvs := createValidEnvironments(t)
		originalSource := createValidTemplateSource(t)
		template, _ := NewDeploymentTemplate(originalSource, originalEnvs, originalSteps)

		retrievedEnvs := template.Environments()
		modifiedEnv, _ := vos.NewEnvironment("MODIFIED", "desc", "val")
		retrievedEnvs[0] = modifiedEnv

		if template.Environments()[0].Name() == "MODIFIED" {
			t.Errorf("El estado interno (environments) fue modificado externamente")
		}
	})

	t.Run("Steps debe devolver una copia", func(t *testing.T) {
		originalSteps := createValidSteps(t)
		originalEnvs := createValidEnvironments(t)
		originalSource := createValidTemplateSource(t)
		template, _ := NewDeploymentTemplate(originalSource, originalEnvs, originalSteps)
		verifications := []vos.VerificationType{vos.VerificationTypeCode}

		retrievedSteps := template.Steps()
		cmd, _ := vos.NewCommandDefinition("c", "c")
		modifiedStep, _ := entities.NewStepDefinition("MODIFIED", verifications, []vos.CommandDefinition{cmd})
		retrievedSteps[0] = modifiedStep

		if template.Steps()[0].Name() == "MODIFIED" {
			t.Errorf("El estado interno (steps) fue modificado externamente")
		}
	})
}

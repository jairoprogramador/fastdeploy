package aggregates

import (
	"testing"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/entities"
	deploymentvos "github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/vos"
	sharedvos "github.com/jairoprogramador/fastdeploy-core/internal/domain/shared/vos"
)

const ENV_STAGING = "staging"
const STEP_TEST = "test"

// --- Helper para crear datos de prueba ---
func createValidTemplateSource(t *testing.T) sharedvos.TemplateSource {
	t.Helper()
	source, err := sharedvos.NewTemplateSource("http://test.com/repo.git", "main")
	if err != nil {
		t.Fatalf("fallo al crear helper TemplateSource: %v", err)
	}
	return source
}

func createValidEnvironments(t *testing.T) []deploymentvos.Environment {
	t.Helper()
	env1, err1 := deploymentvos.NewEnvironment(ENV_STAGING, "stag")
	env2, err2 := deploymentvos.NewEnvironment("production", "prod")
	if err1 != nil || err2 != nil {
		t.Fatalf("fallo al crear helpers Environment: %v, %v", err1, err2)
	}
	return []deploymentvos.Environment{env1, env2}
}

func createValidSteps(t *testing.T) []entities.StepDefinition {
	t.Helper()
	validVariable, _ := deploymentvos.NewVariable("test-var", "hello")
	verifications := []deploymentvos.Trigger{deploymentvos.ScopeCode}
	verifications2 := []deploymentvos.Trigger{deploymentvos.ScopeVars}
	cmd, err := deploymentvos.NewCommandDefinition("test-cmd", "echo")
	if err != nil {
		t.Fatalf("fallo al crear helper CommandDefinition: %v", err)
	}
	step1, err1 := entities.NewStepDefinition(STEP_TEST, verifications, []deploymentvos.CommandDefinition{cmd}, []deploymentvos.Variable{validVariable})
	step2, err2 := entities.NewStepDefinition("deploy", verifications2, []deploymentvos.CommandDefinition{cmd}, []deploymentvos.Variable{validVariable})
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

	verifications := []deploymentvos.Trigger{deploymentvos.ScopeCode}
	validVariable, _ := deploymentvos.NewVariable("test-var", "hello")

	// Crear un environment duplicado para el caso de prueba de fallo
	envDupe, _ := deploymentvos.NewEnvironment(ENV_STAGING, "stag-dupe")
	dupeEnvs := append(validEnvs, envDupe)

	// Crear un paso duplicado
	cmd, _ := deploymentvos.NewCommandDefinition("cmd", "c")
	stepDupe, _ := entities.NewStepDefinition(STEP_TEST, verifications, []deploymentvos.CommandDefinition{cmd}, []deploymentvos.Variable{validVariable})
	dupeSteps := append(validSteps, stepDupe)

	testCases := []struct {
		testName     string
		source       sharedvos.TemplateSource
		environments []deploymentvos.Environment
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
			environments: []deploymentvos.Environment{},
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
		modifiedEnv, _ := deploymentvos.NewEnvironment("MODIFIED", "val")
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
		verifications := []deploymentvos.Trigger{deploymentvos.ScopeCode}
		validVariable, _ := deploymentvos.NewVariable("test-var", "hello")

		retrievedSteps := template.Steps()
		cmd, _ := deploymentvos.NewCommandDefinition("c", "c")
		modifiedStep, _ := entities.NewStepDefinition("MODIFIED", verifications, []deploymentvos.CommandDefinition{cmd}, []deploymentvos.Variable{validVariable})
		retrievedSteps[0] = modifiedStep

		if template.Steps()[0].Name() == "MODIFIED" {
			t.Errorf("El estado interno (steps) fue modificado externamente")
		}
	})
}

package aggregates

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	deploymentaggregates "github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/aggregates"
	deploymententities "github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/entities"
	deploymentvos "github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/vos"
	orchestrationvos "github.com/jairoprogramador/fastdeploy-core/internal/domain/orchestration/vos"
	sharedvos "github.com/jairoprogramador/fastdeploy-core/internal/domain/shared/vos"
)

const ENV_STAGING_NAME = "staging"
const ENV_STAGING_VALUE = "stag"

const STEP_TEST = "test"
const STEP_SUPPLY = "supply"
const STEP_DEPLOY = "deploy"

const VARIABLE_ENVIRONMENT = "environment"
const VARIABLE_ORDER_ID = "order_id"

// --- Mocks y Helpers ---
type MockResolver struct {
	mock.Mock
}

func (m *MockResolver) ResolveOutput(probe orchestrationvos.Output, text string) (orchestrationvos.Output, bool, error) {
	args := m.Called(probe, text)
	return args.Get(0).(orchestrationvos.Output), args.Bool(1), args.Error(2)
}

func (m *MockResolver) ResolveTemplate(template string, variables map[string]orchestrationvos.Output) (string, error) {
	args := m.Called(template, variables)
	return args.String(0), args.Error(1)
}

func (m *MockResolver) ResolvePath(path string, variables map[string]orchestrationvos.Output) error {
	args := m.Called(path, variables)
	return args.Error(0)
}

func createTestTemplate(t *testing.T) *deploymentaggregates.DeploymentTemplate {
	t.Helper()
	source, _ := sharedvos.NewTemplateSource("http://test.com/repo.git", "main")
	env, _ := deploymentvos.NewEnvironment(ENV_STAGING_NAME, ENV_STAGING_VALUE)
	cmd, _ := deploymentvos.NewCommandDefinition("cmd", "echo")
	verifications := []deploymentvos.Trigger{deploymentvos.ScopeCode}
	validVariable, _ := deploymentvos.NewVariable("test-var", "hello")

	stepTest, _ := deploymententities.NewStepDefinition(STEP_TEST, verifications, []deploymentvos.CommandDefinition{cmd}, []deploymentvos.Variable{validVariable})
	stepSupply, _ := deploymententities.NewStepDefinition(STEP_SUPPLY, verifications, []deploymentvos.CommandDefinition{cmd}, []deploymentvos.Variable{validVariable})
	stepDeploy, _ := deploymententities.NewStepDefinition(STEP_DEPLOY, verifications, []deploymentvos.CommandDefinition{cmd}, []deploymentvos.Variable{validVariable})

	template, err := deploymentaggregates.NewDeploymentTemplate(
		source,
		[]deploymentvos.Environment{env},
		[]deploymententities.StepDefinition{stepTest, stepSupply, stepDeploy})

	assert.NoError(t, err)
	return template
}

func TestNewOrder(t *testing.T) {
	template := createTestTemplate(t)
	env := template.Environments()[0]
	orderID := orchestrationvos.NewOrderID()

	testCases := []struct {
		testName          string
		finalStepName     string
		skippedSteps      map[string]struct{}
		expectError       bool
		expectedStepCount int
		expectedSkipped   map[string]bool
	}{
		{
			testName:          "Creacion valida hasta el final",
			finalStepName:     STEP_DEPLOY,
			skippedSteps:      nil,
			expectError:       false,
			expectedStepCount: 3,
			expectedSkipped:   map[string]bool{STEP_TEST: false, STEP_SUPPLY: false, STEP_DEPLOY: false},
		},
		{
			testName:          "Creacion valida hasta un paso intermedio",
			finalStepName:     STEP_SUPPLY,
			skippedSteps:      nil,
			expectError:       false,
			expectedStepCount: 2,
		},
		{
			testName:          "Creacion valida omitiendo un paso",
			finalStepName:     STEP_DEPLOY,
			skippedSteps:      map[string]struct{}{STEP_SUPPLY: {}},
			expectError:       false,
			expectedStepCount: 3,
			expectedSkipped:   map[string]bool{STEP_TEST: false, STEP_SUPPLY: true, STEP_DEPLOY: false},
		},
		{
			testName:      "Fallo por paso final invalido",
			finalStepName: "non-existent-step",
			skippedSteps:  nil,
			expectError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			order, err := NewOrder(orderID, template, env, tc.finalStepName, tc.skippedSteps, nil)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, order)
				assert.Len(t, order.StepsRecord(), tc.expectedStepCount)
				assert.Equal(t, orchestrationvos.OrderStatusInProgress, order.Status())

				// Verificar el mapa de variables inicial
				assert.Contains(t, order.Outputs(), VARIABLE_ENVIRONMENT)
				assert.Contains(t, order.Outputs(), VARIABLE_ORDER_ID)
				assert.Equal(t, ENV_STAGING_VALUE, order.Outputs()[VARIABLE_ENVIRONMENT])
				assert.Equal(t, orderID.String(), order.Outputs()[VARIABLE_ORDER_ID])

				// Verificar pasos omitidos
				if tc.expectedSkipped != nil {
					for _, stepExec := range order.StepsRecord() {
						shouldBeSkipped := tc.expectedSkipped[stepExec.Name()]
						isSkipped := stepExec.Status() == orchestrationvos.StepStatusSkipped
						assert.Equal(t, shouldBeSkipped, isSkipped, "El estado de omision para el paso '%s' no es el esperado", stepExec.Name())
					}
				}
			}
		})
	}
}

func TestOrder_MarkCommandAsCompleted_StateTransition(t *testing.T) {
	template := createTestTemplate(t)
	env := template.Environments()[0]
	resolver := new(MockResolver)

	t.Run("Transicion a Successful", func(t *testing.T) {
		order, _ := NewOrder(orchestrationvos.NewOrderID(), template, env, STEP_SUPPLY, nil, nil)

		// Simular ejecución exitosa
		err := order.FinalizeCommand(STEP_TEST, "cmd", "resolved", "log", 0, resolver)
		assert.NoError(t, err)
		assert.Equal(t, orchestrationvos.OrderStatusInProgress, order.Status())

		err = order.FinalizeCommand(STEP_SUPPLY, "cmd", "resolved", "log", 0, resolver)
		assert.NoError(t, err)
		assert.Equal(t, orchestrationvos.OrderStatusSuccessful, order.Status()) // Último paso tuvo éxito
	})

	t.Run("Transicion a Failed", func(t *testing.T) {
		order, _ := NewOrder(orchestrationvos.NewOrderID(), template, env, STEP_DEPLOY, nil, nil)

		// Simular ejecución exitosa del primer paso
		err := order.FinalizeCommand(STEP_TEST, "cmd", "resolved", "log", 0, resolver)
		assert.NoError(t, err)
		assert.Equal(t, orchestrationvos.OrderStatusInProgress, order.Status())

		// Simular fallo en el segundo paso
		err = order.FinalizeCommand(STEP_SUPPLY, "cmd", "resolved", "log", 1, resolver)
		assert.NoError(t, err)
		assert.Equal(t, orchestrationvos.OrderStatusFailed, order.Status())
	})

	t.Run("Recopilacion de variables", func(t *testing.T) {
		cmdWithProbeName := "cmd-probe"
		variableNewVarName := "new_var"
		variableNewVarValue := "secret"
		outputTextVariableValue := fmt.Sprintf("val=%s", variableNewVarValue)
		verifications := []deploymentvos.Trigger{deploymentvos.ScopeCode}
		validVariable, _ := deploymentvos.NewVariable("test-var", "hello")

		outputProbe, _ := deploymentvos.NewOutput(variableNewVarName, "val=(.*)")
		cmdWithProbe, _ := deploymentvos.NewCommandDefinition(
			cmdWithProbeName, "echo", deploymentvos.WithOutputs([]deploymentvos.Output{outputProbe}))
		stepWithProbe, _ := deploymententities.NewStepDefinition(STEP_TEST, verifications, []deploymentvos.CommandDefinition{cmdWithProbe}, []deploymentvos.Variable{validVariable})

		templateProbe, _ := deploymentaggregates.NewDeploymentTemplate(
			template.Source(),
			[]deploymentvos.Environment{env},
			[]deploymententities.StepDefinition{stepWithProbe})

		newVar, _ := orchestrationvos.NewOutputFromNameAndValue(variableNewVarName, variableNewVarValue)
		resolver.On("ExtractVariable", outputProbe, outputTextVariableValue).Return(newVar, true, nil).Once()

		order, _ := NewOrder(orchestrationvos.NewOrderID(), templateProbe, env, STEP_TEST, nil, nil)

		err := order.FinalizeCommand(STEP_TEST, cmdWithProbeName, "resolved", outputTextVariableValue, 0, resolver)
		assert.NoError(t, err)
		assert.Contains(t, order.Outputs(), variableNewVarName)
		assert.Equal(t, variableNewVarValue, order.Outputs()[variableNewVarName])
		resolver.AssertExpectations(t)
	})
}

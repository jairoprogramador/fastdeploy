package aggregates

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	deploymentaggregates "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/aggregates"
	deploymententities "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/entities"
	deploymentvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
)

const ENV_STAGING_NAME = "staging"
const ENV_STAGING_VALUE = "stag"

const STEP_TEST = "test"
const STEP_SUPPLY = "supply"
const STEP_DEPLOY = "deploy"

const VARIABLE_ENVIRONMENT = "environment"
const VARIABLE_ORDER_ID = "order_id"

// --- Mocks y Helpers ---
type MockVariableResolver struct {
	mock.Mock
}

func (m *MockVariableResolver) ExtractVariable(probe deploymentvos.OutputProbe, text string) (vos.Variable, bool, error) {
	args := m.Called(probe, text)
	return args.Get(0).(vos.Variable), args.Bool(1), args.Error(2)
}

func (m *MockVariableResolver) Interpolate(template string, variables map[string]vos.Variable) (string, error) {
	args := m.Called(template, variables)
	return args.String(0), args.Error(1)
}

func (m *MockVariableResolver) ProcessTemplateFile(srcPath, destPath string, variables map[string]vos.Variable) error {
	args := m.Called(srcPath, destPath, variables)
	return args.Error(0)
}

func createTestTemplate(t *testing.T) *deploymentaggregates.DeploymentTemplate {
	t.Helper()
	source, _ := deploymentvos.NewTemplateSource("http://test.com/repo.git", "main")
	env, _ := deploymentvos.NewEnvironment(ENV_STAGING_NAME, "staging description", ENV_STAGING_VALUE)
	cmd, _ := deploymentvos.NewCommandDefinition("cmd", "echo")
	verifications := []deploymentvos.VerificationType{deploymentvos.VerificationTypeCode}
	stepTest, _ := deploymententities.NewStepDefinition(STEP_TEST, verifications, []deploymentvos.CommandDefinition{cmd})
	stepSupply, _ := deploymententities.NewStepDefinition(STEP_SUPPLY, verifications, []deploymentvos.CommandDefinition{cmd})
	stepDeploy, _ := deploymententities.NewStepDefinition(STEP_DEPLOY, verifications, []deploymentvos.CommandDefinition{cmd})

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
	orderID := vos.NewOrderID()

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
				assert.Len(t, order.StepExecutions(), tc.expectedStepCount)
				assert.Equal(t, vos.OrderStatusInProgress, order.Status())

				// Verificar el mapa de variables inicial
				assert.Contains(t, order.VariableMap(), VARIABLE_ENVIRONMENT)
				assert.Contains(t, order.VariableMap(), VARIABLE_ORDER_ID)
				assert.Equal(t, ENV_STAGING_VALUE, order.VariableMap()[VARIABLE_ENVIRONMENT].Value())
				assert.Equal(t, orderID.String(), order.VariableMap()[VARIABLE_ORDER_ID].Value())

				// Verificar pasos omitidos
				if tc.expectedSkipped != nil {
					for _, stepExec := range order.StepExecutions() {
						shouldBeSkipped := tc.expectedSkipped[stepExec.Name()]
						isSkipped := stepExec.Status() == vos.StepStatusSkipped
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
	resolver := new(MockVariableResolver)

	t.Run("Transicion a Successful", func(t *testing.T) {
		order, _ := NewOrder(vos.NewOrderID(), template, env, STEP_SUPPLY, nil, nil)

		// Simular ejecución exitosa
		err := order.MarkCommandAsCompleted(STEP_TEST, "cmd", "resolved", "log", 0, resolver)
		assert.NoError(t, err)
		assert.Equal(t, vos.OrderStatusInProgress, order.Status())

		err = order.MarkCommandAsCompleted(STEP_SUPPLY, "cmd", "resolved", "log", 0, resolver)
		assert.NoError(t, err)
		assert.Equal(t, vos.OrderStatusSuccessful, order.Status()) // Último paso tuvo éxito
	})

	t.Run("Transicion a Failed", func(t *testing.T) {
		order, _ := NewOrder(vos.NewOrderID(), template, env, STEP_DEPLOY, nil, nil)

		// Simular ejecución exitosa del primer paso
		err := order.MarkCommandAsCompleted(STEP_TEST, "cmd", "resolved", "log", 0, resolver)
		assert.NoError(t, err)
		assert.Equal(t, vos.OrderStatusInProgress, order.Status())

		// Simular fallo en el segundo paso
		err = order.MarkCommandAsCompleted(STEP_SUPPLY, "cmd", "resolved", "log", 1, resolver)
		assert.NoError(t, err)
		assert.Equal(t, vos.OrderStatusFailed, order.Status())
	})

	t.Run("Recopilacion de variables", func(t *testing.T) {
		cmdWithProbeName := "cmd-probe"
		variableNewVarName := "new_var"
		variableNewVarValue := "secret"
		outputTextVariableValue := fmt.Sprintf("val=%s", variableNewVarValue)
		verifications := []deploymentvos.VerificationType{deploymentvos.VerificationTypeCode}

		outputProbe, _ := deploymentvos.NewOutputProbe(variableNewVarName, "description", "val=(.*)")
		cmdWithProbe, _ := deploymentvos.NewCommandDefinition(
			cmdWithProbeName, "echo", deploymentvos.WithOutputs([]deploymentvos.OutputProbe{outputProbe}))
		stepWithProbe, _ := deploymententities.NewStepDefinition(STEP_TEST, verifications, []deploymentvos.CommandDefinition{cmdWithProbe})

		templateProbe, _ := deploymentaggregates.NewDeploymentTemplate(
			template.Source(),
			[]deploymentvos.Environment{env},
			[]deploymententities.StepDefinition{stepWithProbe})

		newVar, _ := vos.NewVariable(variableNewVarName, variableNewVarValue)
		resolver.On("ExtractVariable", outputProbe, outputTextVariableValue).Return(newVar, true, nil).Once()

		order, _ := NewOrder(vos.NewOrderID(), templateProbe, env, STEP_TEST, nil, nil)

		err := order.MarkCommandAsCompleted(STEP_TEST, cmdWithProbeName, "resolved", outputTextVariableValue, 0, resolver)
		assert.NoError(t, err)
		assert.Contains(t, order.VariableMap(), variableNewVarName)
		assert.Equal(t, variableNewVarValue, order.VariableMap()[variableNewVarName].Value())
		resolver.AssertExpectations(t)
	})
}

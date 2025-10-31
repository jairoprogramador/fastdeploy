package entities

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	depvos "github.com/jairoprogramador/fastdeploy-core/internal/domain/template/vos"
	orcvos "github.com/jairoprogramador/fastdeploy-core/internal/domain/executor/vos"
)

const CMD_NAME = "test-cmd"

// --- Helpers de prueba ---
func createTestCommandDef(t *testing.T, opts ...depvos.CommandOption) depvos.CommandDefinition {
	t.Helper()
	def, err := depvos.NewCommandDefinition(CMD_NAME, "echo 'hello'", opts...)
	if err != nil {
		t.Fatalf("fallo al crear helper CommandDefinition: %v", err)
	}
	return def
}

func TestNewCommandExecution(t *testing.T) {
	def := createTestCommandDef(t)
	exec := NewCommandRecord(def)

	assert.NotNil(t, exec)
	assert.Equal(t, CMD_NAME, exec.Name())
	assert.Equal(t, orcvos.CommandStatusPending, exec.Status())
}

func TestCommandExecution_Execute(t *testing.T) {
	// --- Definiciones de sondas para los tests ---
	probe_success := `SUCCESS`
	probe_failed := `FAILED`
	log_probe_success := fmt.Sprintf("Command finished: %s", probe_success)
	log_probe_failed := fmt.Sprintf("Command finished: %s", probe_failed)
	log_probe_extract := `app version="123"`
	methodName_extractVariable := "ExtractVariable"

	probeExtract, _ := depvos.NewOutput("version", `version="(\d+)"`)
	probeValidate, _ := depvos.NewOutput("", probe_success)
	outputExtract, _ := depvos.NewOutput("version", "123")
	varExtract := orcvos.NewOutput(outputExtract)

	testCases := []struct {
		testName                string
		def                     depvos.CommandDefinition
		exitCode                int
		log                     string
		setupMock               func(*MockResolver)
		expectedStatus          orcvos.CommandStatus
		expectedOutputVarsCount int
		expectError             bool
	}{
		{
			testName:       "Fallo por exit code no cero",
			def:            createTestCommandDef(t),
			exitCode:       1,
			log:            "command failed",
			setupMock:      func(m *MockResolver) {}, // El mock no debería ser llamado
			expectedStatus: orcvos.CommandStatusFailed,
		},
		{
			testName:       "Exito sin sondas de salida",
			def:            createTestCommandDef(t),
			exitCode:       0,
			log:            "OK",
			setupMock:      func(m *MockResolver) {},
			expectedStatus: orcvos.CommandStatusSuccessful,
		},
		{
			testName: "Exito con sonda de validacion que coincide",
			def:      createTestCommandDef(t, depvos.WithOutputs([]depvos.Output{probeValidate})),
			exitCode: 0,
			log:      log_probe_success,
			setupMock: func(m *MockResolver) {
				m.On(methodName_extractVariable, probeValidate, log_probe_success).Return(orcvos.Output{}, true, nil)
			},
			expectedStatus: orcvos.CommandStatusSuccessful,
		},
		{
			testName: "Fallo con sonda de validacion que no coincide",
			def:      createTestCommandDef(t, depvos.WithOutputs([]depvos.Output{probeValidate})),
			exitCode: 0,
			log:      log_probe_failed,
			setupMock: func(m *MockResolver) {
				m.On(methodName_extractVariable, probeValidate, log_probe_failed).Return(orcvos.Output{}, false, nil)
			},
			expectedStatus: orcvos.CommandStatusFailed,
		},
		{
			testName: "Exito con sonda de extraccion que coincide",
			def:      createTestCommandDef(t, depvos.WithOutputs([]depvos.Output{probeExtract})),
			exitCode: 0,
			log:      log_probe_extract,
			setupMock: func(m *MockResolver) {
				m.On(methodName_extractVariable, probeExtract, log_probe_extract).Return(varExtract, true, nil)
			},
			expectedStatus:          orcvos.CommandStatusSuccessful,
			expectedOutputVarsCount: 1,
		},
		{
			testName: "Fallo por error del extractor",
			def:      createTestCommandDef(t, depvos.WithOutputs([]depvos.Output{probeValidate})),
			exitCode: 0,
			log:      "some log",
			setupMock: func(m *MockResolver) {
				m.On(methodName_extractVariable, mock.Anything, mock.Anything).Return(orcvos.Output{}, false, errors.New("boom"))
			},
			expectedStatus: orcvos.CommandStatusFailed,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			exec := NewCommandRecord(tc.def)
			resolver := new(MockResolver)
			tc.setupMock(resolver)

			err := exec.Finalize("resolved cmd", tc.log, tc.exitCode, resolver)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedStatus, exec.Status())
			assert.Len(t, exec.Outputs(), tc.expectedOutputVarsCount)
			resolver.AssertExpectations(t) // Verifica que el mock fue llamado como se esperaba
		})
	}
}

func TestCommandExecution_CannotExecuteTwice(t *testing.T) {
	commandExec := NewCommandRecord(createTestCommandDef(t))
	resolver := new(MockResolver)

	// Primera ejecución (exitosa)
	err := commandExec.Finalize("cmd", "log", 0, resolver)
	assert.NoError(t, err)
	assert.Equal(t, orcvos.CommandStatusSuccessful, commandExec.Status())

	// Segunda ejecución (debería fallar)
	err = commandExec.Finalize("cmd2", "log2", 0, resolver)
	assert.Error(t, err)
}

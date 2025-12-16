package aggregates_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/entities"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockStepCommandExecutor struct {
	mock.Mock
}

func (m *MockStepCommandExecutor) Execute(ctx context.Context, command vos.Command, currentVars vos.VariableSet, workspaceRoot string) *vos.ExecutionResult {
	args := m.Called(ctx, command, currentVars, workspaceRoot)
	return args.Get(0).(*vos.ExecutionResult)
}

func TestStepExecutor_Execute_Success(t *testing.T) {
	cmdExecutor := new(MockStepCommandExecutor)
	stepExecutor := aggregates.NewStepExecutor(cmdExecutor)

	varName1, varValue1 := "var1", "val1"
	varInitName1, varInitValue1 := "init", "true"
	pathRoot := "/root"
	log1, log2 := "log1", "log2"

	output1, _ := vos.NewCommandOutput(varName1, varValue1)
	cmd1, _ := vos.NewCommand("cmd1", "echo "+log1, vos.WithOutputs([]vos.CommandOutput{output1}))
	cmd2, _ := vos.NewCommand("cmd2", "echo "+log2)

	step, _ := entities.NewStep("test-step",
		entities.WithCommands([]vos.Command{cmd1, cmd2}),
		entities.WithWorkspaceRoot(pathRoot),
	)
	initialVars := vos.VariableSet{varInitName1: varInitValue1}

	// Mockea la primera llamada
	cmdExecutor.On("Execute", mock.Anything, cmd1, initialVars, pathRoot).Return(&vos.ExecutionResult{
		Status:     vos.Success,
		Logs:       log1,
		OutputVars: vos.VariableSet{varName1: varValue1},
	}).Once()

	// La segunda llamada debe recibir las variables actualizadas
	expectedVarsForCmd2 := vos.VariableSet{varInitName1: varInitValue1, varName1: varValue1}
	cmdExecutor.On("Execute", mock.Anything, cmd2, expectedVarsForCmd2, pathRoot).Return(&vos.ExecutionResult{
		Status: vos.Success,
		Logs:   log2,
	}).Once()

	// Act
	result, err := stepExecutor.Execute(context.Background(), &step, initialVars)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, vos.Success, result.Status)
	assert.Contains(t, result.Logs, log1)
	assert.Contains(t, result.Logs, log2)
	assert.Equal(t, vos.VariableSet{varInitName1: varInitValue1, varName1: varValue1}, result.OutputVars)
	assert.NoError(t, result.Error)
	cmdExecutor.AssertExpectations(t)
}

func TestStepExecutor_Execute_StopsOnFailure(t *testing.T) {
	cmdExecutor := new(MockStepCommandExecutor)
	stepExecutor := aggregates.NewStepExecutor(cmdExecutor)

	cmd1, _ := vos.NewCommand("cmd1", "failing command")
	cmd2, _ := vos.NewCommand("cmd2", "should not run")
	step, _ := entities.NewStep("fail-step", entities.WithCommands([]vos.Command{cmd1, cmd2}))
	failError := errors.New("command failed")

	// Mockea la primera llamada para que falle
	cmdExecutor.On("Execute", mock.Anything, cmd1, mock.Anything, mock.Anything).Return(&vos.ExecutionResult{
		Status: vos.Failure,
		Error:  failError,
		Logs:   "error log",
	}).Once()

	// Act
	result, err := stepExecutor.Execute(context.Background(), &step, vos.VariableSet{})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, vos.Failure, result.Status)
	assert.ErrorIs(t, result.Error, failError)
	assert.Contains(t, result.Logs, "error log")
	// Verifica que el mock de cmdExecutor solo fue llamado una vez, probando que el bucle se detuvo
	cmdExecutor.AssertExpectations(t)
	cmdExecutor.AssertNumberOfCalls(t, "Execute", 1)
}

func TestStepExecutor_Execute_StopsOnIrrecoverableError(t *testing.T) {
	// Arrange
	cmdExecutor := new(MockStepCommandExecutor)
	stepExecutor := aggregates.NewStepExecutor(cmdExecutor)

	cmd1, _ := vos.NewCommand("cmd1", "failing command")
	step, _ := entities.NewStep("fail-step", entities.WithCommands([]vos.Command{cmd1}))
	irrecoverableError := errors.New("irrecoverable")

	// Mockea la llamada para que devuelva un error en el campo Error del resultado
	cmdExecutor.On("Execute", mock.Anything, cmd1, mock.Anything, mock.Anything).Return(&vos.ExecutionResult{
		Error: irrecoverableError,
	}).Once()

	// Act
	result, err := stepExecutor.Execute(context.Background(), &step, vos.VariableSet{})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, vos.Failure, result.Status)
	assert.ErrorIs(t, result.Error, irrecoverableError)
	// Verifica que el mock fue llamado
	cmdExecutor.AssertExpectations(t)
}

package entities

import (
	"testing"
	"github.com/stretchr/testify/assert"
	deploymententities "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/entities"
	deploymentvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
)

// --- Helpers de prueba ---
func createTestStepDef(t *testing.T, name string, commandNames []string) deploymententities.StepDefinition {
	t.Helper()
	var cmds []deploymentvos.CommandDefinition
	for _, cmdName := range commandNames {
		cmd, err := deploymentvos.NewCommandDefinition(cmdName, "do something")
		assert.NoError(t, err)
		cmds = append(cmds, cmd)
	}
	stepDef, err := deploymententities.NewStepDefinition(name, cmds)
	assert.NoError(t, err)
	return stepDef
}

func TestNewStepExecution(t *testing.T) {
	t.Run("Creacion valida", func(t *testing.T) {
		stepDef := createTestStepDef(t, "test", []string{"cmd1", "cmd2"})
		stepExec, err := NewStepExecution(stepDef)

		assert.NoError(t, err)
		assert.NotNil(t, stepExec)
		assert.Equal(t, "test", stepExec.Name())
		assert.Equal(t, vos.StepStatusPending, stepExec.Status())
		assert.Len(t, stepExec.CommandExecutions(), 2)
	})

	t.Run("Fallo por StepDefinition sin comandos", func(t *testing.T) {
		stepDef, _ := deploymententities.NewStepDefinition("test", []deploymentvos.CommandDefinition{})
		_, err := NewStepExecution(stepDef)
		assert.Error(t, err)
	})

	t.Run("Fallo por StepDefinition sin nombre", func(t *testing.T) {
		stepDef, _ := deploymententities.NewStepDefinition("", []deploymentvos.CommandDefinition{})
		_, err := NewStepExecution(stepDef)
		assert.Error(t, err)
	})
}

func TestStepExecution_StateTransitions(t *testing.T) {
	resolver := new(MockVariableResolver)
	// No necesitamos configurar el mock, ya que la lógica de Execute no es el foco aquí.

	t.Run("Transicion a InProgress y luego a Successful", func(t *testing.T) {
		stepDef := createTestStepDef(t, "test", []string{"cmd1", "cmd2"})
		stepExec, _ := NewStepExecution(stepDef)
		assert.Equal(t, vos.StepStatusPending, stepExec.Status())

		// Completar el primer comando con éxito
		err := stepExec.CompleteCommand("cmd1", "resolved", "log", 0, resolver)
		assert.NoError(t, err)
		assert.Equal(t, vos.StepStatusInProgress, stepExec.Status())

		// Completar el segundo comando con éxito
		err = stepExec.CompleteCommand("cmd2", "resolved", "log", 0, resolver)
		assert.NoError(t, err)
		assert.Equal(t, vos.StepStatusSuccessful, stepExec.Status())
	})

	t.Run("Transicion a Failed comando no encontrado", func(t *testing.T) {
		stepDef := createTestStepDef(t, "test", []string{"cmd1", "cmd2"})
		stepExec, _ := NewStepExecution(stepDef)
		assert.Equal(t, vos.StepStatusPending, stepExec.Status())

		err := stepExec.CompleteCommand("cmd3", "resolved", "log", 0, resolver)
		assert.Error(t, err)
	})

	t.Run("Transicion a Failed en el primer comando", func(t *testing.T) {
		stepDef := createTestStepDef(t, "test", []string{"cmd1", "cmd2"})
		stepExec, _ := NewStepExecution(stepDef)

		// Fallar el primer comando
		err := stepExec.CompleteCommand("cmd1", "resolved", "log", 1, resolver)
		assert.NoError(t, err)
		assert.Equal(t, vos.StepStatusFailed, stepExec.Status())
	})

	t.Run("Transicion a Failed en el segundo comando", func(t *testing.T) {
		stepDef := createTestStepDef(t, "test", []string{"cmd1", "cmd2"})
		stepExec, _ := NewStepExecution(stepDef)

		// Completar el primer comando con éxito
		err := stepExec.CompleteCommand("cmd1", "resolved", "log", 0, resolver)
		assert.NoError(t, err)
		assert.Equal(t, vos.StepStatusInProgress, stepExec.Status())

		// Fallar el segundo comando
		err = stepExec.CompleteCommand("cmd2", "resolved", "log", 1, resolver)
		assert.NoError(t, err)
		assert.Equal(t, vos.StepStatusFailed, stepExec.Status())
	})
}

func TestStepExecution_Skip(t *testing.T) {
	stepDef := createTestStepDef(t, "test", []string{"cmd1"})
	stepExec, _ := NewStepExecution(stepDef)

	stepExec.Skip()
	assert.Equal(t, vos.StepStatusSkipped, stepExec.Status())

	// Un paso ya fallido no puede ser omitido
	stepExec.status = vos.StepStatusFailed
	stepExec.Skip()
	assert.Equal(t, vos.StepStatusFailed, stepExec.Status())
}

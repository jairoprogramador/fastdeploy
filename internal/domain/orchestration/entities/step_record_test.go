package entities

import (
	"testing"

	depEnt "github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/entities"
	depVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/orchestration/vos"
	"github.com/stretchr/testify/assert"
)

// --- Helpers de prueba ---
func createTestStepDef(t *testing.T, name string, commandNames []string) depEnt.StepDefinition {
	t.Helper()
	var cmds []depVos.CommandDefinition
	for _, cmdName := range commandNames {
		cmd, err := depVos.NewCommandDefinition(cmdName, "do something")
		assert.NoError(t, err)
		cmds = append(cmds, cmd)
	}
	verifications := []depVos.Trigger{depVos.ScopeCode}
	validVariable, _ := depVos.NewVariable("test-var", "hello")
	stepDef, err := depEnt.NewStepDefinition(name, verifications, cmds, []depVos.Variable{validVariable})
	assert.NoError(t, err)
	return stepDef
}

func TestNewStepRecord(t *testing.T) {
	t.Run("Creacion valida", func(t *testing.T) {
		stepDef := createTestStepDef(t, "test", []string{"cmd1", "cmd2"})
		stepExec := NewStepRecord(stepDef)

		assert.NotNil(t, stepExec)
		assert.Equal(t, "test", stepExec.Name())
		assert.Equal(t, vos.StepStatusPending, stepExec.Status())
		assert.Len(t, stepExec.Commands(), 2)
	})

	t.Run("Fallo por StepDefinition sin comandos", func(t *testing.T) {
		verifications := []depVos.Trigger{depVos.ScopeCode}
		validVariable, _ := depVos.NewVariable("test-var", "hello")
		stepDef, _ := depEnt.NewStepDefinition("test", verifications, []depVos.CommandDefinition{}, []depVos.Variable{validVariable})
		stepExec := NewStepRecord(stepDef)

		assert.NotNil(t, stepExec)
		assert.Equal(t, "test", stepExec.Name())
		assert.Equal(t, vos.StepStatusPending, stepExec.Status())
		assert.Len(t, stepExec.Commands(), 0)
	})

	t.Run("Fallo por StepDefinition sin nombre", func(t *testing.T) {
		verifications := []depVos.Trigger{depVos.ScopeCode}
		validVariable, _ := depVos.NewVariable("test-var", "hello")
		stepDef, _ := depEnt.NewStepDefinition("", verifications, []depVos.CommandDefinition{}, []depVos.Variable{validVariable})
		stepExec := NewStepRecord(stepDef)

		assert.NotNil(t, stepExec)
		assert.Equal(t, "test", stepExec.Name())
		assert.Equal(t, vos.StepStatusPending, stepExec.Status())
		assert.Len(t, stepExec.Commands(), 0)
	})
}

func TestStepRecord_StateTransitions(t *testing.T) {
	resolver := new(MockResolver)
	// No necesitamos configurar el mock, ya que la lógica de Execute no es el foco aquí.

	t.Run("Transicion a InProgress y luego a Successful", func(t *testing.T) {
		stepDef := createTestStepDef(t, "test", []string{"cmd1", "cmd2"})
		stepExec := NewStepRecord(stepDef)
		assert.Equal(t, vos.StepStatusPending, stepExec.Status())

		// Completar el primer comando con éxito
		err := stepExec.FinalizeCommand("cmd1", "resolved", "log", 0, resolver)
		assert.NoError(t, err)
		assert.Equal(t, vos.StepStatusInProgress, stepExec.Status())

		// Completar el segundo comando con éxito
		err = stepExec.FinalizeCommand("cmd2", "resolved", "log", 0, resolver)
		assert.NoError(t, err)
		assert.Equal(t, vos.StepStatusSuccessful, stepExec.Status())
	})

	t.Run("Transicion a Failed comando no encontrado", func(t *testing.T) {
		stepDef := createTestStepDef(t, "test", []string{"cmd1", "cmd2"})
		stepExec := NewStepRecord(stepDef)
		assert.Equal(t, vos.StepStatusPending, stepExec.Status())

		err := stepExec.FinalizeCommand("cmd3", "resolved", "log", 0, resolver)
		assert.Error(t, err)
	})

	t.Run("Transicion a Failed en el primer comando", func(t *testing.T) {
		stepDef := createTestStepDef(t, "test", []string{"cmd1", "cmd2"})
		stepExec := NewStepRecord(stepDef)

		// Fallar el primer comando
		err := stepExec.FinalizeCommand("cmd1", "resolved", "log", 1, resolver)
		assert.NoError(t, err)
		assert.Equal(t, vos.StepStatusFailed, stepExec.Status())
	})

	t.Run("Transicion a Failed en el segundo comando", func(t *testing.T) {
		stepDef := createTestStepDef(t, "test", []string{"cmd1", "cmd2"})
		stepExec := NewStepRecord(stepDef)

		// Completar el primer comando con éxito
		err := stepExec.FinalizeCommand("cmd1", "resolved", "log", 0, resolver)
		assert.NoError(t, err)
		assert.Equal(t, vos.StepStatusInProgress, stepExec.Status())

		// Fallar el segundo comando
		err = stepExec.FinalizeCommand("cmd2", "resolved", "log", 1, resolver)
		assert.NoError(t, err)
		assert.Equal(t, vos.StepStatusFailed, stepExec.Status())
	})
}

func TestStepRecord_Skip(t *testing.T) {
	stepDef := createTestStepDef(t, "test", []string{"cmd1"})
	stepExec := NewStepRecord(stepDef)

	stepExec.Skip()
	assert.Equal(t, vos.StepStatusSkipped, stepExec.Status())

	// Un paso ya fallido no puede ser omitido
	stepExec.status = vos.StepStatusFailed
	stepExec.Skip()
	assert.Equal(t, vos.StepStatusFailed, stepExec.Status())
}

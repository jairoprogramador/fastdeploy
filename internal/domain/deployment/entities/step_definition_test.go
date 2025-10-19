package entities

import (
	"reflect"
	"testing"

	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
)

func TestNewStepDefinition(t *testing.T) {
	validCmd, _ := vos.NewCommandDefinition("test-cmd", "echo 'hello'")
	validTriggers := []vos.Trigger{vos.ScopeCode}
	validVariable, _ := vos.NewVariable("test-var", "hello")

	testCases := []struct {
		testName      string
		name          string
		commands      []vos.CommandDefinition
		verifications []vos.Trigger
		variables     []vos.Variable
		expectError   bool
	}{
		{
			testName:      "Creacion valida",
			name:          "test",
			commands:      []vos.CommandDefinition{validCmd},
			verifications: validTriggers,
			variables:     []vos.Variable{validVariable},
			expectError:   false,
		},
		{
			testName:      "Fallo por nombre vacio",
			name:          "",
			commands:      []vos.CommandDefinition{validCmd},
			verifications: validTriggers,
			variables:     []vos.Variable{validVariable},
			expectError:   true,
		},
		{
			testName:      "Fallo por lista de comandos vacia",
			name:          "test",
			commands:      []vos.CommandDefinition{},
			verifications: validTriggers,
			variables:     []vos.Variable{validVariable},
			expectError:   true,
		},
		{
			testName:      "Fallo por lista de comandos nula",
			name:          "test",
			commands:      nil,
			verifications: validTriggers,
			variables:     []vos.Variable{validVariable},
			expectError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			stepDef, err := NewStepDefinition(tc.name, tc.verifications, tc.commands, tc.variables)

			if tc.expectError {
				if err == nil {
					t.Errorf("Se esperaba un error, pero no se obtuvo ninguno")
				}
			} else {
				if err != nil {
					t.Errorf("No se esperaba un error, pero se obtuvo: %v", err)
				}
				if stepDef.Name() != tc.name {
					t.Errorf("Se esperaba el nombre '%s', pero se obtuvo '%s'", tc.name, stepDef.Name())
				}
				if !reflect.DeepEqual(stepDef.TriggersInt(), tc.verifications) {
					t.Errorf("La lista de verificaciones no coincide")
				}
				if !reflect.DeepEqual(stepDef.Commands(), tc.commands) {
					t.Errorf("La lista de comandos no coincide")
				}
			}
		})
	}
}

func TestStepDefinition_DefensiveCopying(t *testing.T) {
	validVariable, _ := vos.NewVariable("test-var", "hello")
	validCmd1, _ := vos.NewCommandDefinition("cmd1", "echo 1")
	validCmd2, _ := vos.NewCommandDefinition("cmd2", "echo 2")
	originalCommands := []vos.CommandDefinition{validCmd1, validCmd2}
	originalTriggers := []vos.Trigger{vos.ScopeCode}

	stepDef, _ := NewStepDefinition("test-step", originalTriggers, originalCommands, []vos.Variable{validVariable})

	retrievedCommands := stepDef.Commands()
	if len(retrievedCommands) == 0 {
		t.Fatal("Commands() no debería devolver un slice vacío")
	}

	retrievedTriggers := stepDef.TriggersInt()
	if len(retrievedTriggers) == 0 {
		t.Fatal("Triggers() no debería devolver un slice vacío")
	}

	// Modificamos el slice obtenido
	modifiedCmd, _ := vos.NewCommandDefinition("MODIFIED", "echo 'modified'")
	retrievedCommands[0] = modifiedCmd

	// Verificamos que el estado interno no haya cambiado
	internalCommands := stepDef.Commands()
	if internalCommands[0].Name() != "cmd1" {
		t.Errorf("El estado interno de StepDefinition fue modificado. Se esperaba 'cmd1', se obtuvo '%s'", internalCommands[0].Name())
	}

	retrievedTriggers[0] = int(vos.ScopeRecipe)

	if stepDef.TriggersInt()[0] != int(vos.ScopeRecipe) {
		t.Errorf("El estado interno de StepDefinition fue modificado. Se esperaba 'code', se obtuvo '%d'", stepDef.TriggersInt()[0])
	}
}

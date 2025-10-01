package entities

import (
	"reflect"
	"testing"

	"github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
)

func TestNewStepDefinition(t *testing.T) {
	// Helper para crear un comando válido reutilizable
	validCmd, _ := vos.NewCommandDefinition("test-cmd", "echo 'hello'")
	validVerifications := []vos.VerificationType{vos.VerificationTypeCode}

	testCases := []struct {
		testName    string
		name        string
		commands    []vos.CommandDefinition
		verifications []vos.VerificationType
		expectError bool
	}{
		{
			testName:    "Creacion valida",
			name:        "test",
			commands:    []vos.CommandDefinition{validCmd},
			verifications: validVerifications,
			expectError: false,
		},
		{
			testName:    "Fallo por nombre vacio",
			name:        "",
			commands:    []vos.CommandDefinition{validCmd},
			verifications: validVerifications,
			expectError: true,
		},
		{
			testName:    "Fallo por lista de comandos vacia",
			name:        "test",
			commands:    []vos.CommandDefinition{},
			verifications: validVerifications,
			expectError: true,
		},
		{
			testName:    "Fallo por lista de comandos nula",
			name:        "test",
			commands:    nil,
			verifications: validVerifications,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			stepDef, err := NewStepDefinition(tc.name, tc.verifications, tc.commands)

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
				if !reflect.DeepEqual(stepDef.VerificationTypes(), tc.verifications) {
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
	validCmd1, _ := vos.NewCommandDefinition("cmd1", "echo 1")
	validCmd2, _ := vos.NewCommandDefinition("cmd2", "echo 2")
	originalCommands := []vos.CommandDefinition{validCmd1, validCmd2}
	originalVerifications := []vos.VerificationType{vos.VerificationTypeCode}

	stepDef, _ := NewStepDefinition("test-step", originalVerifications, originalCommands)

	retrievedCommands := stepDef.Commands()
	if len(retrievedCommands) == 0 {
		t.Fatal("Commands() no debería devolver un slice vacío")
	}

	retrievedVerifications := stepDef.VerificationTypes()
	if len(retrievedVerifications) == 0 {
		t.Fatal("VerificationTypes() no debería devolver un slice vacío")
	}

	// Modificamos el slice obtenido
	modifiedCmd, _ := vos.NewCommandDefinition("MODIFIED", "echo 'modified'")
	retrievedCommands[0] = modifiedCmd

	// Verificamos que el estado interno no haya cambiado
	internalCommands := stepDef.Commands()
	if internalCommands[0].Name() != "cmd1" {
		t.Errorf("El estado interno de StepDefinition fue modificado. Se esperaba 'cmd1', se obtuvo '%s'", internalCommands[0].Name())
	}

	retrievedVerifications[0] = vos.VerificationTypeEnv

	if stepDef.VerificationTypes()[0] != vos.VerificationTypeEnv {
		t.Errorf("El estado interno de StepDefinition fue modificado. Se esperaba 'code', se obtuvo '%s'", stepDef.VerificationTypes()[0])
	}
}

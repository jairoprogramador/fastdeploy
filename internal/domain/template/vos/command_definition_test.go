package vos

import (
	"reflect"
	"testing"
)

func TestNewCommandDefinition(t *testing.T) {
	// Creamos un OutputProbe válido para usar en los tests
	validProbe, _ := NewOutput("test_probe", ".*")

	testCases := []struct {
		testName       string
		name           string
		cmd    string
		opts           []CommandOption
		expectError    bool
		expectedCmdDef CommandDefinition
	}{
		{
			testName:    "Creacion basica valida",
			name:        "run-script",
			cmd: "bash script.sh",
			opts:        nil,
			expectError: false,
			expectedCmdDef: CommandDefinition{
				name: "run-script",
				cmd:  "bash script.sh",
			},
		},
		{
			testName:    "Creacion valida con todas las opciones",
			name:        "full-command",
			cmd: "kubectl apply -f .",
			opts: []CommandOption{
				WithWorkdir("./k8s"),
				WithTemplateFiles([]string{"a.yaml", "b.yaml"}),
				WithOutputs([]Output{validProbe}),
			},
			expectError: false,
			expectedCmdDef: CommandDefinition{
				name:          "full-command",
				cmd:           "kubectl apply -f .",
				workdir:       "./k8s",
				templateFiles: []string{"a.yaml", "b.yaml"},
				outputs:       []Output{validProbe},
			},
		},
		{
			testName:    "Fallo por nombre vacio",
			name:        "",
			cmd: "some-cmd",
			opts:        nil,
			expectError: true,
		},
		{
			testName:    "Fallo por plantilla de comando vacia",
			name:        "no-cmd",
			cmd: "",
			opts:        nil,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			cmdDef, err := NewCommandDefinition(tc.name, tc.cmd, tc.opts...)

			if tc.expectError {
				if err == nil {
					t.Errorf("Se esperaba un error, pero no se obtuvo ninguno")
				}
			} else {
				if err != nil {
					t.Errorf("No se esperaba un error, pero se obtuvo: %v", err)
				}
				// reflect.DeepEqual es útil para comparar structs complejos.
				if !reflect.DeepEqual(cmdDef, tc.expectedCmdDef) {
					t.Errorf("La CommandDefinition creada no coincide con la esperada.\nObtenido: %+v\nEsperado: %+v", cmdDef, tc.expectedCmdDef)
				}
			}
		})
	}
}

func TestCommandDefinition_DefensiveCopying(t *testing.T) {
	t.Run("TemplateFiles debe devolver una copia", func(t *testing.T) {
		originalFiles := []string{"a.yaml", "b.yaml"}
		cmdDef, _ := NewCommandDefinition("cmd", "do", WithTemplateFiles(originalFiles))

		// Obtenemos el slice del getter
		retrievedFiles := cmdDef.TemplateFiles()
		if len(retrievedFiles) == 0 {
			t.Fatal("TemplateFiles() no debería devolver un slice vacío")
		}

		// Modificamos el slice obtenido
		retrievedFiles[0] = "MODIFIED.yaml"

		// Verificamos que el estado interno del objeto no haya cambiado
		if cmdDef.TemplateFiles()[0] != "a.yaml" {
			t.Errorf("El estado interno de CommandDefinition fue modificado. Se esperaba 'a.yaml', se obtuvo '%s'", cmdDef.TemplateFiles()[0])
		}
	})

	t.Run("Outputs debe devolver una copia", func(t *testing.T) {
		probe, _ := NewOutput("name", ".*")
		originalOutputs := []Output{probe}
		cmdDef, _ := NewCommandDefinition("cmd", "do", WithOutputs(originalOutputs))

		retrievedOutputs := cmdDef.Outputs()
		if len(retrievedOutputs) == 0 {
			t.Fatal("Outputs() no debería devolver un slice vacío")
		}

		// Modificamos el slice obtenido
		retrievedOutputs[0].name = "MODIFIED"

		// Verificamos que el estado interno no haya cambiado
		if cmdDef.Outputs()[0].name != "name" {
			t.Errorf("El estado interno de CommandDefinition fue modificado. Se esperaba 'name', se obtuvo '%s'", cmdDef.Outputs()[0].name)
		}
	})
}

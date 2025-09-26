package main

import (
	"fmt"
	"log"
	"strings"

	values "github.com/jairoprogramador/fastdeploy/internal/domain/step/values"
	"github.com/spf13/cobra"
)

func NewTestCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test [environment]",
		Short: "Ejecuta las pruebas de calidad del software en un entorno específico.",
		Long: `Este comando ejecuta pruebas unitarias, de integración y otros análisis para asegurar la calidad del código.
Si no se especifica un entorno, se usará 'local' por defecto.`,
		Aliases: []string{"t"},
	}

	var validEnvironments []string

	environments := GetEnvironmentRepository()

	for _, env := range environments {
		validEnvironments = append(validEnvironments, env) // Guardamos los nombres para el mensaje de error
		envCmd := &cobra.Command{
			Use:   env,
			Short: fmt.Sprintf("Ejecuta las pruebas para el entorno %s", env),
			Run: func(cmd *cobra.Command, args []string) {
				runTestForEnvironment(cmd, env)
			},
		}
		cmd.AddCommand(envCmd)
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return runTestForEnvironment(cmd, "local")
		}

		if args[0] == "local" {
			return runTestForEnvironment(cmd, "local")
		}

		invalidEnv := args[0]
		errMsg := fmt.Sprintf("el entorno '%s' no es válido", invalidEnv)

		if len(validEnvironments) > 0 {
			suggestion := fmt.Sprintf("Los entornos disponibles son: %s", strings.Join(validEnvironments, ", "))
			return fmt.Errorf("%s. %s", errMsg, suggestion)
		}

		return fmt.Errorf("%s. No se encontraron entornos configurados", errMsg)
	}

	return cmd
}

func runTestForEnvironment(cmd *cobra.Command, environment string) error {
	log.Printf("Iniciando pruebas para el entorno: %s\n", environment)

	if err := GetCommandExecutor().ExecuteCommand(environment, values.StepTest, []string{}); err != nil {
		log.Fatalf("Error: %v", err)
		return err
	}
	return nil
}

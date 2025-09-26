package main

import (
	"fmt"
	"log"
	"strings"

	values "github.com/jairoprogramador/fastdeploy/internal/domain/step/values"
	"github.com/spf13/cobra"
)

func NewPackageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "package [environment]",
		Short:   "Ejecuta el empaquetado de la aplicación para un entorno específico.",
		Long:    `Este comando ejecuta el empaquetado de la aplicación. Si no se especifica un entorno, se usará 'local' por defecto.`,
		Aliases: []string{"p"},
	}

	var validEnvironments []string

	environments := GetEnvironmentRepository()

	for _, env := range environments {
		validEnvironments = append(validEnvironments, env)
		envCmd := &cobra.Command{
			Use:   env,
			Short: fmt.Sprintf("Ejecuta el empaquetado para el entorno %s", env),
			Run: func(cmd *cobra.Command, args []string) {
				runPackageForEnvironment(cmd, env)
			},
		}
		cmd.AddCommand(envCmd)
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			runPackageForEnvironment(cmd, "local")
			return nil
		}

		if args[0] == "local" {
			runPackageForEnvironment(cmd, "local")
			return nil
		}

		invalidEnv := args[0]
		errMsg := fmt.Sprintf("el entorno '%s' no es válido", invalidEnv)

		if len(validEnvironments) > 0 {
			suggestion := fmt.Sprintf("Los entornos disponibles son: %s", strings.Join(validEnvironments, ", "))
			return fmt.Errorf("%s. %s", errMsg, suggestion)
		}

		return fmt.Errorf("%s. No se encontraron entornos configurados", errMsg)
	}

	AddSkipFlags(cmd, getSkipStepsPackage())
	return cmd
}

func runPackageForEnvironment(cmd *cobra.Command, environment string) error {
	log.Printf("Iniciando empaquetado para el entorno: %s\n", environment)

	if err := GetCommandExecutor().ExecuteCommand(environment, values.StepPackage, []string{}); err != nil {
		log.Fatalf("Error: %v", err)
		return err
	}
	return nil
}

func getSkipStepsPackage() []string {
	return []string{values.StepTest, values.StepSupply}
}

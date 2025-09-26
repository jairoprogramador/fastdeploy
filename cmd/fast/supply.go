package main

import (
	"fmt"
	"log"
	"strings"

	values "github.com/jairoprogramador/fastdeploy/internal/domain/step/values"
	"github.com/spf13/cobra"
)

func NewSupplyCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:     "supply [environment]",
		Short:   "Ejecuta el suministro de la aplicación en un entorno específico.",
		Long:    `Este comando ejecuta el suministro de la aplicación. Si no se especifica un entorno, se usará 'local' por defecto.`,
		Aliases: []string{"s"},
	}

	var validEnvironments []string

	environments := GetEnvironmentRepository()

	for _, env := range environments {
		validEnvironments = append(validEnvironments, env)
		envCmd := &cobra.Command{
			Use:   env,
			Short: fmt.Sprintf("Ejecuta el suministro para el entorno %s", env),
			Run: func(cmd *cobra.Command, args []string) {
				runSupplyForEnvironment(cmd, env)
			},
		}
		cmd.AddCommand(envCmd)
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return runSupplyForEnvironment(cmd, "local")
		}

		if args[0] == "local" {
			return runSupplyForEnvironment(cmd, "local")
		}

		invalidEnv := args[0]
		errMsg := fmt.Sprintf("el entorno '%s' no es válido", invalidEnv)

		if len(validEnvironments) > 0 {
			suggestion := fmt.Sprintf("Los entornos disponibles son: %s", strings.Join(validEnvironments, ", "))
			return fmt.Errorf("%s. %s", errMsg, suggestion)
		}

		return fmt.Errorf("%s. No se encontraron entornos configurados", errMsg)
	}

	AddSkipFlags(cmd, getSkipStepsSupply())
	return cmd
}

func runSupplyForEnvironment(cmd *cobra.Command, environment string) error {
	log.Printf("Iniciando suministro para el entorno: %s\n", environment)

	if err := GetCommandExecutor().ExecuteCommand(environment, values.StepSupply, []string{}); err != nil {
		log.Fatalf("Error: %v", err)
		return err
	}
	return nil
}

func getSkipStepsSupply() []string {
	return []string{values.StepTest}
}

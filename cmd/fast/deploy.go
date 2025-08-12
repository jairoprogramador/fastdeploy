package main

import (
	"fmt"
	"log"

	"github.com/jairoprogramador/fastdeploy/internal/adapters/cli"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/cli/utils"
	"github.com/jairoprogramador/fastdeploy/internal/constants"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/commands"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/project"
	"github.com/spf13/cobra"
)

func NewDeployCmd() *cobra.Command {
	skippableSteps := []string{constants.StepTest, constants.StepSupply}

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Ejecuta el despliegue de la aplicación.",
		Long:  `Este comando ejecuta el despliegue de la aplicación.`,
		Run: func(cmd *cobra.Command, args []string) {
			projectTechnology := "java" // o "node"

			projectEntity, err := project.Load()
			if err != nil {
				log.Fatalf("Error al leer datos del proyecto: %v", err)
			}

			repositoryFilePath, err := utils.GetRepositoryPath(*projectEntity)
			if err != nil {
				log.Fatalf("Error al obtener ruta del repositorio: %v", err)
			}

			factory, err := cli.GetStrategyFactory(projectTechnology, repositoryFilePath)
			if err != nil {
				log.Fatalf("Error al obtener la fábrica de estrategias: %v", err)
			}

			allCommands := map[string]commands.Command{
				constants.StepTest:    commands.NewTestCommand(factory.CreateTestStrategy()),
				constants.StepSupply:  commands.NewSupplyCommand(factory.CreateSupplyStrategy()),
				constants.StepPackage: commands.NewPackageCommand(factory.CreatePackageStrategy()),
				constants.StepDeploy:  commands.NewDeployCommand(factory.CreateDeployStrategy()),
			}

			skipFlags := utils.GetSkipFlags(cmd, skippableSteps)

			executionOrder := []string{constants.StepTest, constants.StepSupply, constants.StepPackage, constants.StepDeploy}

			firstCommand, err := utils.BuildDynamicChain(allCommands, skipFlags, executionOrder)
			if err != nil {
				log.Fatalf("Error al construir la cadena de comandos: %v", err)
			}

			if firstCommand != nil {
				pipelineContext := context.NewPipelineContext()
				pipelineContext.Set(constants.Technology, projectEntity.Technology)

				if err := firstCommand.Execute(pipelineContext); err != nil {
					log.Fatalf("Error al ejecutar el comando: %v", err)
				}
			} else {
				fmt.Println("No se seleccionaron pasos para ejecutar. Saliendo...")
			}
		},
	}
	utils.AddSkipFlags(cmd, skippableSteps)
	return cmd
}

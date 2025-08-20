package main

import (
	"fmt"
	//"log"

	//factory "github.com/jairoprogramador/fastdeploy/internal/adapters/factory/impl"
	//"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies/manager"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/utils"
	"github.com/jairoprogramador/fastdeploy/internal/constants"
	//"github.com/jairoprogramador/fastdeploy/internal/core/domain/commands"
	//"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/spf13/cobra"
)

func NewSupplyCmd() *cobra.Command {
	skippableSteps := []string{constants.StepTest}

	cmd := &cobra.Command{
		Use:   "supply",
		Short: "Ejecuta el suministro de la aplicación.",
		Long:  `Este comando ejecuta el suministro de la aplicación.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Supply command")
			/* projectTechnology := "java" // o "node"

			projectEntity, err := factory.NewServiceFactory().CreateProjectService().Load()
			if err != nil {
				log.Fatalf("Error al leer datos del proyecto: %v", err)
			}

			repositoryPath, err := factory.NewPathFactory().CreateGitPathResolver().GetDirectoryPath(projectEntity.Repository)
			if err != nil {
				log.Fatalf("Error al obtener ruta del repositorio: %v", err)
			}

			factory, err := manager.NewFactoryManager().GetFactory(projectTechnology, repositoryPath)
			if err != nil {
				log.Fatalf("Error al obtener la fábrica de estrategias: %v", err)
			}

			allCommands := map[string]commands.Command{
				constants.StepTest:   commands.NewTestCommand(factory.CreateTestStrategy()),
				constants.StepSupply: commands.NewSupplyCommand(factory.CreateSupplyStrategy()),
			}

			skipFlags := utils.GetSkipFlags(cmd, skippableSteps)

			executionOrder := []string{constants.StepTest, constants.StepSupply}

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
			} */
		},
	}
	utils.AddSkipFlags(cmd, skippableSteps)
	return cmd
}

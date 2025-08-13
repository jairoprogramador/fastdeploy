package main

import (
	factory "github.com/jairoprogramador/fastdeploy/internal/adapters/factory/impl"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/strategies/manager"
	"github.com/jairoprogramador/fastdeploy/internal/constants"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/commands"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"github.com/spf13/cobra"
	"log"
)

func NewTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Ejecuta las pruebas de calidad del software.",
		Long: `Este comando ejecuta pruebas unitarias, de integración, escaneos de seguridad y otros análisis estáticos
			para asegurar la calidad del código.`,
		Run: func(cmd *cobra.Command, args []string) {
			projectTechnology := "java" // o "node"

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

			testCommand := commands.NewTestCommand(factory.CreateTestStrategy())

			pipelineContext := context.NewPipelineContext()
			pipelineContext.Set(constants.Technology, projectEntity.Technology)

			if err := testCommand.Execute(pipelineContext); err != nil {
				log.Fatalf("Error al ejecutar el comando test: %v", err)
			}
		},
	}
}

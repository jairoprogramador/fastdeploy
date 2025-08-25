package main

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/jairoprogramador/fastdeploy/internal/application/project"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment"
	domainService "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/service"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/project/service"
	app "github.com/jairoprogramador/fastdeploy/internal/application/deployment"
	constantInfra "github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
	constantDomain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/constant"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/factory"
)

func NewTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Ejecuta las pruebas de calidad del software.",
		Long: `Este comando ejecuta pruebas unitarias, de integración, escaneos de seguridad y otros análisis estáticos
			para asegurar la calidad del código.`,
		Run: func(cmd *cobra.Command, args []string) {
			repositoryProject := service.NewFileRepository()
			readerProject := project.NewReader(repositoryProject)

			context := deployment.NewDeploymentContext()
			registryStrategy := factory.NewRegistryStrategy()

			factoryStrategy, err := registryStrategy.Get(constantInfra.FactoryManual)
			if err != nil {
				log.Fatalf("Error al obtener el factory strategy: %v", err)
			}

			commandManager := domainService.NewStepOrchestrator(factoryStrategy)

			executeStep := app.NewExecuteStep(readerProject, context, commandManager)

			if err := executeStep.StartStep(constantDomain.StepTest, []string{}); err != nil {
				log.Fatalf("Error al ejecutar el comando test: %v", err)
			}
		},
	}
}

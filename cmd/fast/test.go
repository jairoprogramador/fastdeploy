package main

import (
	"log"

	app "github.com/jairoprogramador/fastdeploy/internal/application/deployment"
	"github.com/jairoprogramador/fastdeploy/internal/application/project"
	domainContext "github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
	constantDomain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/constant"
	domainService "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/service"
	constantInfra "github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
	contextService "github.com/jairoprogramador/fastdeploy/internal/infrastructure/context/service"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/factory"
	projectService "github.com/jairoprogramador/fastdeploy/internal/infrastructure/project/service"
	"github.com/spf13/cobra"
)

func NewTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Ejecuta las pruebas de calidad del software.",
		Long: `Este comando ejecuta pruebas unitarias, de integración, escaneos de seguridad y otros análisis estáticos
			para asegurar la calidad del código.`,
		Aliases: []string{"t"},
		Run: func(cmd *cobra.Command, args []string) {
			repositoryProject := projectService.NewFileRepository()
			readerProject := project.NewReader(repositoryProject)
			identifier := projectService.NewHashIdentifier()

			context := domainContext.NewDataContext()
			registryStrategy := factory.NewRegistryStrategy()

			factoryStrategy, err := registryStrategy.Get(constantInfra.FactoryManual)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			commandManager := domainService.NewStepOrchestrator(factoryStrategy)

			contextRepository := contextService.NewFileRepository()

			executeStep := app.NewExecuteStep(readerProject, identifier, context, contextRepository, commandManager)

			if err := executeStep.StartStep(constantDomain.StepTest, []string{}); err != nil {
				log.Fatalf("Error: %v", err)
			}
		},
	}
}

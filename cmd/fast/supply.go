package main

import (
	"log"
	projectService "github.com/jairoprogramador/fastdeploy/internal/infrastructure/project/service"
	contextService "github.com/jairoprogramador/fastdeploy/internal/infrastructure/context/service"
	"github.com/jairoprogramador/fastdeploy/internal/application/project"
	domainContext "github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
	app "github.com/jairoprogramador/fastdeploy/internal/application/deployment"
	domainService "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/service"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/factory"
	constantInfra "github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
	constantDomain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/constant"

	"github.com/spf13/cobra"
)

func NewSupplyCmd() *cobra.Command {
	skippableSteps := []string{constantDomain.StepTest}

	cmd := &cobra.Command{
		Use:   "supply",
		Short: "Ejecuta el suministro de la aplicación.",
		Long:  `Este comando ejecuta el suministro de la aplicación.`,
		Aliases: []string{"s"},
		Run: func(cmd *cobra.Command, args []string) {
			repositoryProject := projectService.NewFileRepository()
			readerProject := project.NewReader(repositoryProject)

			context := domainContext.NewDataContext()
			registryStrategy := factory.NewRegistryStrategy()

			factoryStrategy, err := registryStrategy.Get(constantInfra.FactoryManual)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			commandManager := domainService.NewStepOrchestrator(factoryStrategy)

			contextRepository := contextService.NewFileRepository()

			executeStep := app.NewExecuteStep(readerProject, context, contextRepository, commandManager)

			if err := executeStep.StartStep(constantDomain.StepSupply, GetSkipSteps(cmd, skippableSteps)); err != nil {
				log.Fatalf("Error: %v", err)
			}
		},
	}
	AddSkipFlags(cmd, skippableSteps)
	return cmd
}

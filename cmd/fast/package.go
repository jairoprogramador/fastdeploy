package main

import (
	"log"

	constantInfra "github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
	constantDomain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/constant"
	projectService "github.com/jairoprogramador/fastdeploy/internal/infrastructure/project/service"
	contextService "github.com/jairoprogramador/fastdeploy/internal/infrastructure/context/service"
	"github.com/jairoprogramador/fastdeploy/internal/application/project"
	domainContext "github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
	app "github.com/jairoprogramador/fastdeploy/internal/application/deployment"
	domainService "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/service"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/factory"
	"github.com/spf13/cobra"
)

func NewPackageCmd() *cobra.Command {
	skippableSteps := []string{constantDomain.StepTest, constantDomain.StepSupply}

	cmd := &cobra.Command{
		Use:   "package",
		Short: "Ejecuta el empaquetado de la aplicación.",
		Long:  `Este comando ejecuta el empaquetado de la aplicación.`,
		Aliases: []string{"p"},
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

			if err := executeStep.StartStep(constantDomain.StepPackage, GetSkipSteps(cmd, skippableSteps)); err != nil {
				log.Fatalf("Error: %v", err)
			}
		},
	}
	AddSkipFlags(cmd, skippableSteps)
	return cmd
}

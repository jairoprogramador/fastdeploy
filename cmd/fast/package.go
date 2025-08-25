package main

import (
	"log"

	constantInfra "github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
	constantDomain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/constant"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/project/service"
	"github.com/jairoprogramador/fastdeploy/internal/application/project"
	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment"
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
		Run: func(cmd *cobra.Command, args []string) {
			repositoryProject := service.NewFileRepository()
			readerProject := project.NewReader(repositoryProject)

			context := deployment.NewDeploymentContext()
			registryStrategy := factory.NewRegistryStrategy()

			factoryStrategy, err := registryStrategy.Get(constantInfra.FactoryManual)
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			commandManager := domainService.NewStepOrchestrator(factoryStrategy)

			executeStep := app.NewExecuteStep(readerProject, context, commandManager)

			if err := executeStep.StartStep(constantDomain.StepPackage, GetSkipSteps(cmd, skippableSteps)); err != nil {
				log.Fatalf("Error: %v", err)
			}
		},
	}
	AddSkipFlags(cmd, skippableSteps)
	return cmd
}

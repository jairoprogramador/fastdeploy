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

func NewDeployCmd() *cobra.Command {
	skippableSteps := []string{constantDomain.StepTest, constantDomain.StepSupply}

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Ejecuta el despliegue de la aplicación.",
		Long:  `Este comando ejecuta el despliegue de la aplicación.`,
		Run: func(cmd *cobra.Command, args []string) {
			repositoryProject := service.NewFileRepository()
			readerProject := project.NewReader(repositoryProject)

			context := deployment.NewDeploymentContext()
			registryStrategy := factory.NewRegistryStrategy()

			factoryStrategy, err := registryStrategy.Get(constantInfra.FactoryManual)
			if err != nil {
				log.Fatalf("Error al obtener el factory strategy: %v", err)
			}

			commandManager := domainService.NewCommandManager(factoryStrategy)

			executeStep := app.NewExecuteStep(readerProject, context, commandManager)

			if err := executeStep.StartStep(constantDomain.StepDeploy, GetSkipSteps(cmd, skippableSteps)); err != nil {
				log.Fatalf("Error al ejecutar el comando supply: %v", err)
			}
		},
	}
	AddSkipFlags(cmd, skippableSteps)
	return cmd
}

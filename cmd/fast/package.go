package main

import (
	"fmt"
	"log"
	"strings"

	app "github.com/jairoprogramador/fastdeploy/internal/application/deployment"
	"github.com/jairoprogramador/fastdeploy/internal/application/project"
	domainContext "github.com/jairoprogramador/fastdeploy/internal/domain/context/service"
	constantDomain "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/constant"
	domainService "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/service"
	constantInfra "github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
	contextService "github.com/jairoprogramador/fastdeploy/internal/infrastructure/context/service"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/factory"
	deploymentService "github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/service"
	projectService "github.com/jairoprogramador/fastdeploy/internal/infrastructure/project/service"
	"github.com/spf13/cobra"
)

func NewPackageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "package [environment]",
		Short:   "Ejecuta el empaquetado de la aplicación para un entorno específico.",
		Long:    `Este comando ejecuta el empaquetado de la aplicación. Si no se especifica un entorno, se usará 'local' por defecto.`,
		Aliases: []string{"p"},
	}

	repositoryProject := projectService.NewFileRepository()
	readerProject := project.NewReader(repositoryProject)
	proj, err := readerProject.Read()
	var validEnvironments []string
	if err != nil {
		log.Printf("Advertencia: no se ha podido leer el proyecto para crear subcomandos de package: %v", err)
	} else {
		repoName := proj.GetRepository().GetURL().ExtractNameRepository()
		environmentRepository := deploymentService.NewEnvironmentRepository()
		environments, err := environmentRepository.GetEnvironments(repoName)
		if err != nil {
			log.Printf("Advertencia: no se pudieron obtener los entornos para package: %v", err)
		}

		for _, env := range environments {
			envName := env.GetName()
			validEnvironments = append(validEnvironments, envName)
			envCmd := &cobra.Command{
				Use:   envName,
				Short: fmt.Sprintf("Ejecuta el empaquetado para el entorno %s", envName),
				Run: func(cmd *cobra.Command, args []string) {
					runPackageForEnvironment(cmd, envName)
				},
			}
			cmd.AddCommand(envCmd)
		}
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			runPackageForEnvironment(cmd, "local")
			return nil
		}

		if args[0] == "local" {
			runPackageForEnvironment(cmd, "local")
			return nil
		}

		invalidEnv := args[0]
		errMsg := fmt.Sprintf("el entorno '%s' no es válido", invalidEnv)

		if len(validEnvironments) > 0 {
			suggestion := fmt.Sprintf("Los entornos disponibles son: %s", strings.Join(validEnvironments, ", "))
			return fmt.Errorf("%s. %s", errMsg, suggestion)
		}

		return fmt.Errorf("%s. No se encontraron entornos configurados", errMsg)
	}

	AddSkipFlags(cmd, getSkipStepsPackage())
	return cmd
}

func runPackageForEnvironment(cmd *cobra.Command, environment string) {
	log.Printf("Iniciando empaquetado para el entorno: %s\n", environment)

	repositoryProject := projectService.NewFileRepository()
	readerProject := project.NewReader(repositoryProject)
	identifier := projectService.NewHashIdentifier()

	context := domainContext.NewDataContext()
	context.Set(constantInfra.Environment, environment)

	registryStrategy := factory.NewRegistryStrategy()
	factoryStrategy, err := registryStrategy.Get(constantInfra.FactoryManual)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	commandManager := domainService.NewStepOrchestrator(factoryStrategy)

	contextRepository := contextService.NewFileRepository()
	environmentRepository := deploymentService.NewEnvironmentRepository()

	validateEnvironment := domainService.NewValidateEnvironment(environmentRepository)

	executeStep := app.NewExecuteStep(readerProject, identifier, context, contextRepository, commandManager, validateEnvironment)

	if err := executeStep.StartStep(constantDomain.StepPackage, GetSkipSteps(cmd, getSkipStepsPackage())); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func getSkipStepsPackage() []string {
	return []string{constantDomain.StepTest, constantDomain.StepSupply}
}

package main

import (
	"fmt"
	"log"
	appConfig "github.com/jairoprogramador/fastdeploy/internal/application/configuration"
	appProject "github.com/jairoprogramador/fastdeploy/internal/application/project"
	domainFactoryProject "github.com/jairoprogramador/fastdeploy/internal/domain/project/factories"
	domainServiceProject "github.com/jairoprogramador/fastdeploy/internal/domain/project/services"
	domainServiceConfig "github.com/jairoprogramador/fastdeploy/internal/domain/configuration/services"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/configuration"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/project"
	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Inicializa un nuevo proyecto con la configuración de fastDeploy.",
		Long: `Este comando crea el archivo fastDeploy.yaml en el directorio actual 
	con las configuraciones por defecto, como el nombre del proyecto, ID y versión.`,
		Run: func(cmd *cobra.Command, args []string) {

			projectRepository := project.NewFileRepository()
			projectValidator := domainServiceProject.NewValidatorProject()

			configRepository := configuration.NewFileRepository()
			validatorConfig := domainServiceConfig.NewValidatorConfiguration()
			readerConfig := appConfig.NewReader(configRepository, validatorConfig)

			readerProject := appProject.NewReader(projectRepository, projectValidator)
			writerProject := appProject.NewWriter(projectRepository)

			projectFactory := domainFactoryProject.NewProjectFactory()
			projectGit := project.NewGit()
			projectIdentifier := project.NewHashIdentifier()
			projectName := project.NewProjectName()

			projectInitializer := appProject.NewInitializer(readerConfig, readerProject, writerProject, projectFactory, projectGit, projectIdentifier, projectName, projectValidator)
			
			project, err := projectInitializer.Initialize()
			if err != nil {
				log.Fatalf("Error al inicializar el proyecto: %v", err)
			}

			fmt.Printf("Proyecto '%s' inicializado correctamente.\n", project.GetName().Value())
			//fmt.Printf("Archivo de configuración '%s' creado.\n", "fastDeploy.yaml")
		},
	}
}

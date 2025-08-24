package main

import (
	"fmt"
	"log"
	appConfig "github.com/jairoprogramador/fastdeploy/internal/application/configuration/service"
	appProject "github.com/jairoprogramador/fastdeploy/internal/application/project"
	domainFactoryProject "github.com/jairoprogramador/fastdeploy/internal/domain/project/factory"
	infConfigService "github.com/jairoprogramador/fastdeploy/internal/infrastructure/configuration/service"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/project/service"
	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Inicializa un nuevo proyecto con la configuración de fastDeploy.",
		Long: `Este comando crea el archivo fastDeploy.yaml en el directorio actual 
	con las configuraciones por defecto, como el nombre del proyecto, ID y versión.`,
		Run: func(cmd *cobra.Command, args []string) {

			projectRepository := service.NewFileRepository()
			configRepository := infConfigService.NewFileRepository()

			readerConfig := appConfig.NewReader(configRepository)

			readerProject := appProject.NewReader(projectRepository)
			writerProject := appProject.NewWriter(projectRepository)

			projectFactory := domainFactoryProject.NewProjectFactory()
			projectGit := service.NewGitManager()
			projectIdentifier := service.NewHashIdentifier()
			projectName := service.NewProjectName()

			projectInitializer := appProject.NewInitializer(readerConfig, readerProject, writerProject, projectFactory, projectGit, projectIdentifier, projectName)
			
			project, err := projectInitializer.Initialize()
			if err != nil {
				log.Fatalf("Error al inicializar el proyecto: %v", err)
			}

			fmt.Printf("Proyecto '%s' inicializado correctamente.\n", project.GetName().Value())
		},
	}
}

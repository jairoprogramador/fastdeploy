package main

import (
	"fmt"
	"os"

	appConfig "github.com/jairoprogramador/fastdeploy/internal/application/configuration"
	appProject "github.com/jairoprogramador/fastdeploy/internal/application/project"
	domainFactoryProject "github.com/jairoprogramador/fastdeploy/internal/domain/project/factories"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/configuration"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/project"
	"github.com/spf13/cobra"
)

var rootCmd *cobra.Command

func NewRootCmd() *cobra.Command {
	if rootCmd != nil {
		return rootCmd
	}

	rootCmd = &cobra.Command {
		Use:   "fd",
		Short: "CLI para gestionar despliegues de aplicaciones",
		Long:  `Una herramienta de línea de comandos para gestionar el despliegue de aplicaciones en diferentes ambientes con dependencias configurables.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Name() == "init" || cmd.Name() == "config" {
				return
			}

			projectRepository := project.NewFileRepository()
			configRepository := configuration.NewFileRepository()

			readerConfig := appConfig.NewReader(configRepository)

			readerProject := appProject.NewReader(projectRepository)
			writerProject := appProject.NewWriter(projectRepository)

			projectFactory := domainFactoryProject.NewProjectFactory()
			projectGit := project.NewGit()
			projectIdentifier := project.NewHashIdentifier()
			projectName := project.NewProjectName()

			projectInitializer := appProject.NewInitializer(readerConfig, readerProject, writerProject, projectFactory, projectGit, projectIdentifier, projectName)

			isInitialized, err := projectInitializer.IsInitialized()
			if err != nil {
				fmt.Println("Error al verificar si el proyecto está inicializado:", err)
				os.Exit(1)
			}

			if !isInitialized {
				fmt.Println("El despliegue del proyecto no ha sido inicializado.")
				fmt.Println("Por favor, ejecuta 'init' para comenzar.")
				os.Exit(1)
			}

			/* if !factory.NewInitializeFactory().CreateInitialize(configRepository, projectRepository).IsInitialized() {
				fmt.Println("El despliegue del proyecto no ha sido inicializado.")
				fmt.Println("Por favor, ejecuta 'init' para comenzar.")
				os.Exit(1)
			} */
		},
	}

	return rootCmd
}

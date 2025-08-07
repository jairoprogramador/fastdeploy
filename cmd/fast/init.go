package main

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/project"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

func NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Inicializa un nuevo proyecto con la configuración de fastDeploy.",
		Long: `Este comando crea el archivo fastDeploy.yaml en el directorio actual 
	con las configuraciones por defecto, como el nombre del proyecto, ID y versión.`,
		Run: func(cmd *cobra.Command, args []string) {
			initializer := project.NewInitializer()

			if initializer.CheckIfAlreadyInitialized() {
				fmt.Println("¡El proyecto ya ha sido inicializado!")
				return
			}

			projectName, err := getProjectName()
			if err != nil {
				log.Fatalf("Error: %v", err)
			}

			cfg, err := initializer.InitializeProject(projectName)
			if err != nil {
				log.Fatalf("Error al inicializar el proyecto: %v", err)
			}

			fmt.Printf("Proyecto '%s' inicializado correctamente.\n", cfg.ProjectName)
			fmt.Printf("Archivo de configuración '%s' creado.\n", "fastDeploy.yaml")
		},
	}
}

func getProjectName() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio de trabajo: %w", err)
	}
	return filepath.Base(dir), nil
}

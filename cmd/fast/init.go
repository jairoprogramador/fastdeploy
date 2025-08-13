package main

import (
	"fmt"
	"log"

	factory "github.com/jairoprogramador/fastdeploy/internal/adapters/factory/impl"
	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Inicializa un nuevo proyecto con la configuración de fastDeploy.",
		Long: `Este comando crea el archivo fastDeploy.yaml en el directorio actual 
	con las configuraciones por defecto, como el nombre del proyecto, ID y versión.`,
		Run: func(cmd *cobra.Command, args []string) {
			initializer := factory.NewInitializeFactory().CreateInitialize()

			if initializer.IsInitialized() {
				fmt.Println("¡El proyecto ya ha sido inicializado!")
				return
			}

			cfg, err := initializer.Initialize()
			if err != nil {
				log.Fatalf("Error al inicializar el proyecto: %v", err)
			}

			fmt.Printf("Proyecto '%s' inicializado correctamente.\n", cfg.ProjectName)
			fmt.Printf("Archivo de configuración '%s' creado.\n", "fastDeploy.yaml")
		},
	}
}

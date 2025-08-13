package main

import (
	"fmt"
	"log"

	"github.com/jairoprogramador/fastdeploy/internal/adapters/config"
	"github.com/spf13/cobra"
)

var (
	organizationConfig string
	teamNameConfig     string
	repositoryConfig   string
	listConfig         bool
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configura los valores globales de la herramienta.",
		Long:  `Este comando te permite establecer valores por defecto para la organización, el equipo y el repositorio, que se usarán al inicializar nuevos proyectos.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Crear las dependencias
			configService := config.NewConfigFactory().CreateService()

			cfg, err := configService.Load()
			if err != nil {
				log.Fatalf("Error al cargar la configuración existente: %v", err)
			}

			if listConfig {
				fmt.Println("Configuración global de FastDeploy:")
				fmt.Printf("  Organización: %s\n", cfg.Organization)
				fmt.Printf("  Nombre del Equipo: %s\n", cfg.TeamName)
				fmt.Printf("  Repositorio: %s\n", cfg.Repository)
				return
			}

			if organizationConfig != "" || teamNameConfig != "" || repositoryConfig != "" {
				if organizationConfig != "" {
					cfg.Organization = organizationConfig
				}
				if teamNameConfig != "" {
					cfg.TeamName = teamNameConfig
				}
				if repositoryConfig != "" {
					cfg.Repository = repositoryConfig
				}

				if err := configService.Save(*cfg); err != nil {
					log.Fatalf("Error al guardar la configuración: %v", err)
				}

				fmt.Println("Configuración global guardada correctamente.")
				return
			}

			cmd.Help()
		},
	}

	cmd.Flags().StringVarP(&organizationConfig, "organization", "o", "", "Nombre de la organización.")
	cmd.Flags().StringVarP(&teamNameConfig, "teamname", "t", "", "Nombre del equipo.")
	cmd.Flags().StringVarP(&repositoryConfig, "repository", "r", "", "URL del repositorio git por defecto.")
	cmd.Flags().BoolVarP(&listConfig, "list", "l", false, "Muestra la configuración global actual.")

	return cmd
}

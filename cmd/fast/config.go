package main

import (
	"fmt"
	"log"

	"github.com/jairoprogramador/fastdeploy/internal/application/configuration/dto"
	appService "github.com/jairoprogramador/fastdeploy/internal/application/configuration/service"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/configuration/service"
	"github.com/spf13/cobra"
)

var (
	organizationValue  string
	teamNameValue      string
	repositoryUrlValue string
	technologyName     string
	listConfig         bool
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configura los valores globales de la herramienta.",
		Long:  `Este comando te permite establecer valores por defecto para la organización, el equipo y el repositorio, que se usarán al inicializar nuevos proyectos.`,
		Run: func(cmd *cobra.Command, args []string) {

			configRepository := service.NewFileRepository()

			if listConfig {
				readerConfig := appService.NewReader(configRepository)
				showConfig := appService.NewDataShow(readerConfig)

				_, err := showConfig.Show()
				if err != nil {
					log.Fatalf("Error al mostrar la configuración: %v", err)
				}
				return
			} else {
				writerConfig := appService.NewWriter(configRepository)

				configDto := dto.ConfigDto{
					NameOrganization: organizationValue,
					Team:             teamNameValue,
					UrlRepository:    repositoryUrlValue,
					Technology:       technologyName,
				}

				err := writerConfig.Write(configDto)
				if err != nil {
					log.Fatalf("Error al guardar la configuración: %v", err)
				}

				fmt.Println("Configuración guardada correctamente.")
				return
			}
		},
	}

	cmd.Flags().StringVarP(&organizationValue, "organization", "o", "", "asigna nombre de la organización.")
	cmd.Flags().StringVarP(&teamNameValue, "team", "t", "", "asigna nombre del equipo.")
	cmd.Flags().StringVarP(&repositoryUrlValue, "repository", "r", "", "asigna URL del repositorio git por defecto.")
	cmd.Flags().StringVarP(&technologyName, "technology", "n", "", "asigna nombre de la tecnología.")
	cmd.Flags().BoolVarP(&listConfig, "list", "l", false, "Muestra la configuración actual.")

	return cmd
}

package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"github.com/spf13/cobra"

	applic "github.com/jairoprogramador/fastdeploy-core/internal/application"
	appDto "github.com/jairoprogramador/fastdeploy-core/internal/application/dto"
	iAppli "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/application"

	iDom "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/dom"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Inicializa a nuevo proyecto",
	Long:  `Inicializa a nuevo proyecto creando el archivo .fastdeploy/dom.yaml`,
	Run:   runInit,
}

func runInit(cmd *cobra.Command, args []string) {

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error al obtener el directorio de trabajo:", err)
		os.Exit(1)
	}

	userInput := iAppli.NewUserInputProvider()
	idGenerator := iDom.NewShaGenerator()
	domRepository := iDom.NewDomYAMLRepository(workingDir)

	initService := applic.NewInitService(userInput, idGenerator, domRepository)

	req := appDto.InitRequest{
		Ctx:              context.Background(),
		SkipPrompt:       skipPrompt,
		WorkingDirectory: filepath.Base(workingDir),
	}

	if _, err := initService.Run(req); err != nil {
		fmt.Printf("\n❌ Error durante la inicialización: %v\n", err)
		os.Exit(1)
	}
}
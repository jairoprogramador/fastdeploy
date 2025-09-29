package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jairoprogramador/fastdeploy/newinternal/application"
	"github.com/jairoprogramador/fastdeploy/newinternal/application/dto"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/git"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/shell"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/state"
	varsinfra "github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/vars"
)

var (
	// Inyectado en tiempo de compilación
	version = "dev"

	// Paths de configuración
	fastdeployHome string
	reposPath      string
	projsPath      string

	// Flags para el comando de ejecución
	skipTest   bool
	skipSupply bool
)

// rootCmd representa el comando base cuando se llama sin subcomandos
var rootCmd = &cobra.Command{
	Use:     "fd <paso> [ambiente]",
	Short:   "fastdeploy es una herramienta CLI para automatizar despliegues.",
	Long:    `Una herramienta para orquestar despliegues de software a través de diferentes ambientes, utilizando plantillas Git.`,
	Version: version,
	Args:    cobra.RangeArgs(1, 2), // Acepta 1 o 2 argumentos
	Run:     runExecution,
}

func init() {
	// Manejo de la configuración
	cobra.OnInitialize(initConfig)

	// Definir flags para el comando raíz
	rootCmd.Flags().BoolVarP(&skipTest, "skip-test", "t", false, "Omitir el paso 'test'")
	rootCmd.Flags().BoolVarP(&skipSupply, "skip-supply", "s", false, "Omitir el paso 'supply'")
}

func initConfig() {
	// Usar Viper para leer la variable de entorno FASTDEPLOY_HOME
	viper.SetEnvPrefix("FASTDEPLOY")
	viper.AutomaticEnv()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error al obtener el directorio home:", err)
		os.Exit(1)
	}

	// Establecer el valor por defecto si la variable de entorno no está definida
	defaultHome := filepath.Join(homeDir, ".fastdeploy")
	fastdeployHome = viper.GetString("HOME")
	if fastdeployHome == "" {
		fastdeployHome = defaultHome
	}

	reposPath = filepath.Join(fastdeployHome, "repositories")
	projsPath = filepath.Join(fastdeployHome, "projects")
}

// runExecution es la función que se ejecuta para el comando principal.
func runExecution(cmd *cobra.Command, args []string) {
	finalStep := args[0]
	environment := "" // Por defecto, un string vacío
	if len(args) == 2 {
		environment = args[1]
	}

	// --- Ensamblaje de Dependencias ---
	cmdExecutor := shell.NewExecutor()
	varResolver := varsinfra.NewResolver()
	templateRepo := git.NewTemplateRepository(reposPath, cmdExecutor)
	orderRepo, err := state.NewFileOrderRepository(projsPath)
	if err != nil {
		fmt.Println("Error al crear order repository:", err)
		os.Exit(1)
	}
	stepVarRepo := git.NewVariableRepository(reposPath)

	// Crear el servicio de aplicación con todas sus dependencias
	orchestrationService := application.NewOrchestrationService(
		templateRepo,
		orderRepo,
		stepVarRepo,
		varResolver,
		cmdExecutor,
	)

	// --- Preparación de la Solicitud ---
	skippedSteps := make(map[string]struct{})
	if skipTest {
		skippedSteps["test"] = struct{}{}
	}
	if skipSupply {
		skippedSteps["supply"] = struct{}{}
	}

	// TODO: Leer deploy.yaml para obtener la URL del repo y otras variables iniciales
	templateSource, err := vos.NewTemplateSource("https://github.com/jairoprogramador/mydeploy.git", "main")
	if err != nil {
		fmt.Println("error al crear el template source:", err)
		os.Exit(1)
	}

	req := dto.ExecuteOrderRequest{
		Ctx:              context.Background(),
		TemplateSource:   templateSource,
		EnvironmentName:  environment, // Pasamos el ambiente (o un string vacío)
		FinalStepName:    finalStep,
		ProjectName:      "projectName",
		SkippedStepNames: skippedSteps,
	}

	// --- Ejecución del Caso de Uso ---
	order, err := orchestrationService.ExecuteOrder(req)
	if err != nil {
		fmt.Printf("\n❌ %v\n", err)
		// Aunque haya un error, el 'order' puede tener estado parcial que queramos inspeccionar.
		if order != nil {
			fmt.Printf("Estado final de la orden (ID: %s): %s\n", order.ID().String(), order.Status())
		}
		os.Exit(1)
	}

	// La ejecución terminó (exitosa o fallida, pero manejada por el servicio).
	fmt.Printf("\n✅ Orden finalizada (ID: %s) con estado: %s\n", order.ID().String(), order.Status())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

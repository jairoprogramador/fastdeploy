package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jairoprogramador/fastdeploy/newinternal/application"
	"github.com/jairoprogramador/fastdeploy/newinternal/application/dto"
	deploymentvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	orchestrationvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
	domaggregates "github.com/jairoprogramador/fastdeploy/newinternal/domain/dom/aggregates"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/console"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/dom"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/executionstate"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/fingerprint"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/git"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/hasher"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/shell"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/state"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/vars"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/workspace"
)

var (
	version        = "dev"
	fastdeployHome string
	reposPath      string
	projsPath      string
	statePath      string

	// Flags para el comando de ejecución
	skipTest   bool
	skipSupply bool

	// Flag para el comando init
	skipPrompt bool
)

var rootCmd = &cobra.Command{
	Use:     "fd [paso] [ambiente]",
	Short:   "fastdeploy es una herramienta CLI para automatizar despliegues.",
	Long:    `Una herramienta para orquestar despliegues de software a través de diferentes ambientes`,
	Version: version,
	// Validamos los argumentos aquí, pero la lógica de ejecución principal está en RunE
	Args: func(cmd *cobra.Command, args []string) error {
		// Permitir que los comandos sin argumentos (como 'fd --version') pasen
		if len(args) == 0 {
			return nil
		}
		// Si hay argumentos, deben ser 1 o 2.
		if len(args) < 1 || len(args) > 2 {
			return errors.New("se requiere un paso y opcionalmente un ambiente")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Si no se pasaron argumentos, pero se ejecutó RunE,
		// significa que el usuario escribió solo "fd". Mostramos la ayuda.
		if len(args) == 0 {
			return cmd.Help()
		}
		// Si se pasaron argumentos, llamamos a la lógica de ejecución.
		runExecution(cmd, args)
		return nil
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Inicializa un proyecto creando el archivo .fastdeploy/dom.yaml",
	Run:   runInit,
}

func init() {
	cobra.OnInitialize(initConfig)
	// Los flags ahora pertenecen al rootCmd, ya que es el que maneja la ejecución.
	rootCmd.Flags().BoolVarP(&skipTest, "skip-test", "t", false, "Omitir el paso 'test'")
	rootCmd.Flags().BoolVarP(&skipSupply, "skip-supply", "s", false, "Omitir el paso 'supply'")

	initCmd.Flags().BoolVarP(&skipPrompt, "yes", "y", false, "Omitir preguntas y usar valores por defecto")

	rootCmd.AddCommand(initCmd)
}

func initConfig() {
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
	statePath = filepath.Join(fastdeployHome, "state")
}

func runInit(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error al obtener el directorio de trabajo:", err)
		os.Exit(1)
	}

	// --- Ensamblaje de Dependencias para Init ---
	userInput := console.NewUserInputProvider()
	idGen := hasher.NewIDGenerator()
	domRepo, err := dom.NewDomYAMLRepository(workingDir)
	if err != nil {
		fmt.Println("Error al inicializar el repositorio DOM:", err)
		os.Exit(1)
	}

	initService := application.NewInitService(userInput, idGen, domRepo)

	// --- Ejecución del Caso de Uso ---
	req := dto.InitRequest{
		Ctx:              ctx,
		SkipPrompt:       skipPrompt,
		WorkingDirectory: filepath.Base(workingDir),
	}

	if _, err := initService.InitializeDOM(req); err != nil {
		fmt.Printf("\n❌ Error durante la inicialización: %v\n", err)
		os.Exit(1)
	}
}

func runExecution(_ *cobra.Command, args []string) {
	ctx := context.Background()
	finalStep := args[0]
	environment := ""
	if len(args) == 2 {
		environment = args[1]
	}

	workingDir, _ := os.Getwd()

	// --- Ensamblaje de Dependencias Comunes ---
	cmdExecutor := shell.NewExecutor()
	varResolver := vars.NewResolver()
	fpService := fingerprint.NewService()
	workspaceMgr := workspace.NewManager(projsPath)
	templateRepo := git.NewTemplateRepository(reposPath, cmdExecutor)
	orderRepo, _ := state.NewFileOrderRepository(projsPath)
	historyRepo, _ := executionstate.NewGobHistoryRepository(statePath)
	domRepo, _ := dom.NewDomYAMLRepository(workingDir)
	idGen := hasher.NewIDGenerator()
	userInput := console.NewUserInputProvider()

	// --- Lógica de Carga y Verificación del DOM ---
	domModel, err := domRepo.Load(ctx)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("⚠️ El archivo .fastdeploy/dom.yaml no existe. Por favor, ejecuta 'fd init' primero.")
		} else {
			fmt.Printf("❌ Error al cargar .fastdeploy/dom.yaml: %v\n", err)
		}
		os.Exit(1)
	}

	isModified, err := domModel.VerifyAndUpdateIDs(idGen)
	if err != nil {
		fmt.Printf("❌ Error al verificar la integridad del DOM: %v\n", err)
		os.Exit(1)
	}

	if isModified {
		history, _ := historyRepo.Find(ctx, "supply")
		if history != nil && len(history.Receipts()) > 0 {
			fmt.Println("⚠️ Se han detectado cambios en .fastdeploy/dom.yaml que afectan a la identidad del proyecto.")
			confirmed, err := userInput.Confirm(ctx, "¿Continuar? Esto podría causar cambios en la infraestructura existente.", false)
			if err != nil || !confirmed {
				fmt.Println("Operación cancelada.")
				os.Exit(1)
			}
		}
		if err := domRepo.Save(ctx, domModel); err != nil {
			fmt.Printf("❌ Error al guardar los cambios en .fastdeploy/dom.yaml: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✅ IDs del proyecto actualizados en .fastdeploy/dom.yaml.")
	}

	// --- Inyección en el Servicio de Aplicación ---
	orchestrationService := application.NewOrchestrationService(
		templateRepo,
		orderRepo,
		historyRepo,
		varResolver,
		fpService,
		workspaceMgr,
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

	templateSource, err := deploymentvos.NewTemplateSource(domModel.Template().RepositoryURL(), domModel.Template().Ref())
	if err != nil {
		fmt.Printf("❌ Error al crear el template source: %v\n", err)
		os.Exit(1)
	}

	initVars, err := initVars(domModel, environment)
	if err != nil {
		fmt.Printf("❌ Error al crear las variables iniciales: %v\n", err)
		os.Exit(1)
	}

	req := dto.ExecuteOrderRequest{
		Ctx:              ctx,
		TemplateSource:   templateSource,
		EnvironmentName:  environment,
		FinalStepName:    finalStep,
		ProjectName:      domModel.Project().Name(),
		ProjectRootPath:  workingDir,
		SkippedStepNames: skippedSteps,
		InitialVariables: initVars,
	}

	// --- Ejecución del Caso de Uso ---
	order, err := orchestrationService.ExecuteOrder(req)
	if err != nil {
		fmt.Printf("\n❌ %v\n", err)
		if order != nil {
			fmt.Printf("Estado final (ID: %s): %s\n", order.ID().String(), order.Status())
		}
		os.Exit(1)
	}
	if order != nil && order.Status() != orchestrationvos.OrderStatusSuccessful {
		os.Exit(1)
	}
	fmt.Printf("\n✅ Orden finalizada (ID: %s) con estado: %s\n", order.ID().String(), order.Status())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initVars(domModel *domaggregates.DeploymentObjectModel, environment string) ([]orchestrationvos.Variable, error) {
	initVars := []orchestrationvos.Variable{}

	projectId, err := orchestrationvos.NewVariable("project_id", domModel.Project().IdString()[:8])
	if err != nil {
		return nil, err
	}
	initVars = append(initVars, projectId)
	projectName, err := orchestrationvos.NewVariable("project_name", domModel.Project().Name())
	if err != nil {
		return nil, err
	}
	initVars = append(initVars, projectName)
	projectTeam, err := orchestrationvos.NewVariable("project_team", domModel.Project().Team())
	if err != nil {
		return nil, err
	}
	initVars = append(initVars, projectTeam)
	projectVersion, err := orchestrationvos.NewVariable("project_version", domModel.Project().Version())
	if err != nil {
		return nil, err
	}
	initVars = append(initVars, projectVersion)

	env, err := orchestrationvos.NewVariable("environment", environment)
	if err != nil {
		return nil, err
	}
	initVars = append(initVars, env)

	return initVars, nil
}

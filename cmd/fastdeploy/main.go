package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/jairoprogramador/fastdeploy/internal/application"
	"github.com/jairoprogramador/fastdeploy/internal/application/dto"
	applicationports "github.com/jairoprogramador/fastdeploy/internal/application/ports"
	deploymentaggregates "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/aggregates"
	domaggregates "github.com/jairoprogramador/fastdeploy/internal/domain/dom/aggregates"
	domports "github.com/jairoprogramador/fastdeploy/internal/domain/dom/ports"
	executionstateports "github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/ports"
	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/console"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/dom"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/executionstate"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/fingerprint"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/git"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/hasher"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/shell"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/state"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/vars"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/workspace"
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

	userInput := console.NewUserInputProvider()
	idGenerator := hasher.NewIDGenerator()
	domRepository := dom.NewDomYAMLRepository(workingDir)
	initService := application.NewInitService(userInput, idGenerator, domRepository)

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
	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	skippedSteps := make(map[string]struct{})
	if skipTest {
		skippedSteps["test"] = struct{}{}
	}
	if skipSupply {
		skippedSteps["supply"] = struct{}{}
	}

	domRepository := dom.NewDomYAMLRepository(workingDir)
	domModel, err := loadDom(ctx, domRepository)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	cmdExecutor := shell.NewExecutor()

	templateResponse, err := loadTemplate(
		ctx, cmdExecutor, reposPath, domModel)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	validateOrderResponse, err := validateOrder(
		templateResponse.Template, environment, finalStep)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	revisionProject, err := loadRevisionProject(
		ctx, cmdExecutor, workingDir, domModel.Project().Revision(),
		validateOrderResponse.FinalStep)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	domModel.SetProjectRevision(revisionProject)

	historyRepository, _ := executionstate.NewScopeRepository(statePath, domModel.Project().Name())

	err = updateDOM(
		ctx,
		domRepository,
		historyRepository,
		domModel,
		validateOrderResponse.Environment.Name())
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	varsRepository, err := executionstate.NewVarsRepository(projsPath, domModel.Project().Name())
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	stateRepository, err := executionstate.NewStateRepository(statePath, domModel.Project().Name())
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	orchestrationService := createOrchestrationService(
		historyRepository,
		stateRepository,
		cmdExecutor,
		varsRepository,
		templateResponse.TemplatePath,
	)

	orderRequest := createOrderRequest(
		ctx, templateResponse,
		validateOrderResponse,
		workingDir, domModel, skippedSteps)

	orderResponse, err := orchestrationService.ExecuteOrder(orderRequest)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	if orderResponse != nil && orderResponse.Status() != vos.OrderStatusSuccessful {
		os.Exit(1)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func loadRevisionProject(
	ctx context.Context,
	cmdExecutor applicationports.CommandExecutor,
	workingDir string,
	revisionDefault string,
	finalStep string) (string, error) {

	gitManager := git.NewGitManager(cmdExecutor, workingDir)
	revisionProjectService := application.NewRevisionProjectService(gitManager)
	return revisionProjectService.LoadProjectRevision(ctx, revisionDefault, finalStep)
}

func loadDom(
	ctx context.Context,
	domRepository domports.DOMRepository) (*domaggregates.DeploymentObjectModel, error) {

	loadDOMService := application.NewLoadDOMService(domRepository)
	return loadDOMService.Load(ctx)
}

func loadTemplate(
	ctx context.Context,
	executor applicationports.CommandExecutor,
	repositoryPath string,
	domModel *domaggregates.DeploymentObjectModel) (dto.LoadTemplateResponse, error) {

	templateRepository := git.NewTemplateRepository(repositoryPath, executor)
	loadTemplateService := application.NewLoadTemplateService(templateRepository)
	return loadTemplateService.Load(
		ctx, domModel.Template().RepositoryURL(), domModel.Template().Ref())
}

func validateOrder(
	template *deploymentaggregates.DeploymentTemplate,
	environment string,
	finalStep string) (dto.ValidateOrderResponse, error) {

	validateOrderService := application.NewValidateOrderService()

	validateOrderRequest := dto.ValidateOrderRequest{
		Environment: environment,
		FinalStep:   finalStep,
	}
	return validateOrderService.Validate(template, validateOrderRequest)
}

func updateDOM(
	ctx context.Context,
	domRepository domports.DOMRepository,
	historyRepository executionstateports.ScopeRepository,
	domModel *domaggregates.DeploymentObjectModel,
	environment string) error {

	idGenerator := hasher.NewIDGenerator()
	userInput := console.NewUserInputProvider()
	updateDOMService := application.NewUpdateDOMService(domRepository, historyRepository, idGenerator, userInput)
	return updateDOMService.Update(ctx, domModel, environment)
}

func createOrchestrationService(
	historyRepository executionstateports.ScopeRepository,
	stateRepository executionstateports.StateRepository,
	executor applicationports.CommandExecutor,
	varsRepository executionstateports.VarsRepository,
	templatePath string) *application.OrchestrationService {

	varResolver := vars.NewResolver()
	fpService := fingerprint.NewFingerprintService()
	workspaceMgr := workspace.NewManager(projsPath, reposPath)
	orderRepo := state.NewFileOrderRepository(projsPath)
	stepVariableRepo := git.NewStepVariableRepository(templatePath)

	return application.NewOrchestrationService(
		stepVariableRepo,
		orderRepo,
		historyRepository,
		varResolver,
		fpService,
		workspaceMgr,
		executor,
		varsRepository,
		stateRepository,
	)
}

func createOrderRequest(
	ctx context.Context,
	templateResponse dto.LoadTemplateResponse,
	validateOrderResponse dto.ValidateOrderResponse,
	workingDir string,
	domModel *domaggregates.DeploymentObjectModel,
	skippedSteps map[string]struct{}) dto.OrderRequest {

	return dto.OrderRequest{
		Ctx:              ctx,
		Template:         templateResponse.Template,
		TemplatePath:     templateResponse.TemplatePath,
		RepositoryName:   templateResponse.RepositoryName,
		ProjectDom:       domModel,
		Environment:      validateOrderResponse.Environment,
		FinalStep:        validateOrderResponse.FinalStep,
		ProjectPath:      workingDir,
		SkippedStepNames: skippedSteps,
	}
}

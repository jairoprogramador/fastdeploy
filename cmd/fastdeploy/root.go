package fastdeploy

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	applic "github.com/jairoprogramador/fastdeploy/internal/application"
	appDto "github.com/jairoprogramador/fastdeploy/internal/application/dto"
	appPor "github.com/jairoprogramador/fastdeploy/internal/application/ports"
	iAppli "github.com/jairoprogramador/fastdeploy/internal/infrastructure/application"

	domAgg "github.com/jairoprogramador/fastdeploy/internal/domain/dom/aggregates"
	domPor "github.com/jairoprogramador/fastdeploy/internal/domain/dom/ports"
	iDom "github.com/jairoprogramador/fastdeploy/internal/infrastructure/dom"

	staPor "github.com/jairoprogramador/fastdeploy/internal/domain/state/ports"
	staSer "github.com/jairoprogramador/fastdeploy/internal/domain/state/services"

	iStaRep "github.com/jairoprogramador/fastdeploy/internal/infrastructure/state/repository"
	iStaSer "github.com/jairoprogramador/fastdeploy/internal/infrastructure/state/services"

	depAgg "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/aggregates"
	iDeplo "github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment"

	orcVos "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"
	iOrche "github.com/jairoprogramador/fastdeploy/internal/infrastructure/orchestration"

	shaVos "github.com/jairoprogramador/fastdeploy/internal/domain/shared/vos"
)

var (
	version = "dev"
	//commit  = "none"
	//date    = "unknown"

	repositoriesPath string
	projectsPath     string
	statePath        string

	skipTest   bool
	skipSupply bool

	skipPrompt bool
)

var rootCmd = &cobra.Command{
	Use:   "fd [paso] [ambiente]",
	Short: "fastdeploy es una herramienta CLI para automatizar despliegues.",
	Long:  `Una herramienta para orquestar despliegues de software a través de diferentes ambientes`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 2 {
			return errors.New("se requiere un paso y opcionalmente un ambiente")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		runOrder(cmd, args)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = fmt.Sprintf("fd version: %s\n", version)
	rootCmd.SetVersionTemplate(`{{.Version}}`)

	cobra.OnInitialize(initConfig)

	rootCmd.Flags().BoolVarP(&skipTest, "skip-test", "t", false, "Omitir el paso 'test'")
	rootCmd.Flags().BoolVarP(&skipSupply, "skip-supply", "s", false, "Omitir el paso 'supply'")
	initCmd.Flags().BoolVarP(&skipPrompt, "yes", "y", false, "Omitir preguntas y usar valores por defecto")

	rootCmd.AddCommand(initCmd)
}

func initConfig() {
	fastdeployHome := getFastdeployHome()

	repositoriesPath = filepath.Join(fastdeployHome, "repositories")
	projectsPath = filepath.Join(fastdeployHome, "projects")
	statePath = filepath.Join(fastdeployHome, "state")
}

func getFastdeployHome() string {
	viper.SetEnvPrefix("FASTDEPLOY")
	viper.AutomaticEnv()

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error al obtener el directorio home:", err)
		os.Exit(1)
	}

	defaultHome := filepath.Join(userHomeDir, ".fastdeploy")
	fastdeployHome := viper.GetString("HOME")
	if fastdeployHome == "" {
		fastdeployHome = defaultHome
	}
	return fastdeployHome
}

func runOrder(_ *cobra.Command, args []string) {
	ctx := context.Background()
	finalStep := args[0]
	environment := "sand"
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

	domRepository := iDom.NewDomYAMLRepository(workingDir)
	domModel, err := loadDom(domRepository)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	cmdExecutor := iAppli.NewExecutor()

	templateResponse, err := loadTemplate(
		ctx, cmdExecutor, repositoriesPath, environment, domModel)
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

	stateRepository, _ := iStaRep.NewFileFingerprintRepository(
		statePath,
		domModel.Project().Name(),
		templateResponse.RepositoryName,
		validateOrderResponse.Environment.Value())

	err = updateDom(
		ctx,
		domRepository,
		stateRepository,
		domModel)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	varsRepository, err := iStaRep.NewVarsRepository(
		statePath,
		domModel.Project().Name(),
		templateResponse.RepositoryName,
		validateOrderResponse.Environment.Value())
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	orderRequest := createOrderRequest(
		ctx, templateResponse,
		validateOrderResponse,
		workingDir, domModel, skippedSteps)

	orchestrationService := createOrchestrationService(
		stateRepository,
		cmdExecutor,
		varsRepository,
		orderRequest,
		validateOrderResponse.Environment.Value(),
	)

	orderResponse, err := orchestrationService.Run(orderRequest)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	if orderResponse != nil && orderResponse.Status() != orcVos.OrderStatusSuccessful {
		os.Exit(1)
	}
}

func loadRevisionProject(
	ctx context.Context,
	cmdExecutor appPor.CommandExecutor,
	workingDir string,
	revisionDefault string,
	finalStep string) (string, error) {

	gitManager := iAppli.NewGitManager(cmdExecutor, workingDir)
	revisionProjectService := applic.NewRevisionProjectService(gitManager)
	return revisionProjectService.LoadProjectRevision(ctx, revisionDefault, finalStep)
}

func loadDom(
	domRepository domPor.DomRepository) (*domAgg.DeploymentObjectModel, error) {

	loadDOMService := applic.NewLoadDOMService(domRepository)
	return loadDOMService.Load()
}

func loadTemplate(
	ctx context.Context,
	executor appPor.CommandExecutor,
	repositoryPath string,
	environment string,
	domModel *domAgg.DeploymentObjectModel) (appDto.TemplateResponse, error) {

	templateRepository := iDeplo.NewTemplateRepository(repositoryPath, environment, executor)
	loadTemplateService := applic.NewLoadTemplateService(templateRepository)

	templateSource, err := shaVos.NewTemplateSource(domModel.Template().Url(), domModel.Template().Ref())
	if err != nil {
		return appDto.TemplateResponse{}, err
	}
	return loadTemplateService.Load(
		ctx, templateSource)
}

func validateOrder(
	template *depAgg.DeploymentTemplate,
	environment string,
	finalStep string) (appDto.ValidateOrderResponse, error) {

	validateOrderService := applic.NewValidateOrderService()

	validateOrderRequest := appDto.ValidateOrderRequest{
		Environment: environment,
		FinalStep:   finalStep,
	}
	return validateOrderService.Validate(template, validateOrderRequest)
}

func updateDom(
	ctx context.Context,
	domRepository domPor.DomRepository,
	stateRepository staPor.FingerprintRepository,
	domModel *domAgg.DeploymentObjectModel) error {

	idGenerator := iDom.NewShaGenerator()
	userInput := iAppli.NewUserInputProvider()
	updateDOMService := applic.NewUpdateDOMService(domRepository, stateRepository, idGenerator, userInput)
	return updateDOMService.Update(ctx, domModel)
}

func createOrchestrationService(
	stateRepository staPor.FingerprintRepository,
	executor appPor.CommandExecutor,
	varsRepository staPor.VariablesRepository,
	orderRequest appDto.OrderRequest,
	environment string) *applic.ExecuteOrder {

	varResolver := iOrche.NewGoTemplateResolver()

	fpService, err := iStaSer.NewFingerprintService(environment)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	statePolicyService := staSer.NewFingerprintPolicyService()

	workspaceMgr, _ := iAppli.NewManager(
		projectsPath,
		repositoriesPath,
		orderRequest.ProjectDom.Project().Name(),
		orderRequest.RepositoryName,
		environment)

	orderRepo := iOrche.NewFileOrderRepository(
		projectsPath,
		orderRequest.ProjectDom.Project().Name(),
		orderRequest.RepositoryName)

	return applic.NewExecuteOrder(
		orderRepo,
		varResolver,
		fpService,
		workspaceMgr,
		executor,
		varsRepository,
		stateRepository,
		statePolicyService,
	)
}

func createOrderRequest(
	ctx context.Context,
	templateResponse appDto.TemplateResponse,
	validateOrderResponse appDto.ValidateOrderResponse,
	workingDir string,
	domModel *domAgg.DeploymentObjectModel,
	skippedSteps map[string]struct{}) appDto.OrderRequest {

	return appDto.OrderRequest{
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

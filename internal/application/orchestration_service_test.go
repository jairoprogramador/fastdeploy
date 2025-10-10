package application

import (
	"context"
	"fmt"
	"testing"

	"github.com/jairoprogramador/fastdeploy/internal/application/dto"
	deploymentaggregates "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/aggregates"
	deploymententities "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/entities"
	deploymentvos "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
	domaggregates "github.com/jairoprogramador/fastdeploy/internal/domain/dom/aggregates"
	domvos "github.com/jairoprogramador/fastdeploy/internal/domain/dom/vos"
	executionstateaggregates "github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/aggregates"
	executionstatevos "github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/vos"
	orchestrationaggregates "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/aggregates"
	orchestrationservices "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/services"
	orchestrationvos "github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocks for dependencies

type MockStepVariableRepository struct {
	mock.Mock
}

func (m *MockStepVariableRepository) Load(stepName string) ([]orchestrationvos.Variable, error) {
	args := m.Called(stepName)
	return args.Get(0).([]orchestrationvos.Variable), args.Error(1)
}

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Save(order *orchestrationaggregates.Order, projectName string) error {
	args := m.Called(order, projectName)
	return args.Error(0)
}

type MockScopeRepository struct {
	mock.Mock
}

func (m *MockScopeRepository) FindCodeStateHistory() (*executionstateaggregates.ScopeReceiptHistory, error) {
	args := m.Called()
	return args.Get(0).(*executionstateaggregates.ScopeReceiptHistory), args.Error(1)
}

func (m *MockScopeRepository) SaveCodeStateHistory(history *executionstateaggregates.ScopeReceiptHistory) error {
	args := m.Called(history)
	return args.Error(0)
}

func (m *MockScopeRepository) FindStepStateHistory(stepName string) (*executionstateaggregates.ScopeReceiptHistory, error) {
	args := m.Called(stepName)
	return args.Get(0).(*executionstateaggregates.ScopeReceiptHistory), args.Error(1)
}

func (m *MockScopeRepository) SaveStepStateHistory(history *executionstateaggregates.ScopeReceiptHistory, stepName string) error {
	args := m.Called(history, stepName)
	return args.Error(0)
}

type MockVariableResolver struct {
	mock.Mock
}

func (m *MockVariableResolver) ExtractVariable(probe deploymentvos.OutputProbe, text string) (orchestrationvos.Variable, bool, error) {
	args := m.Called(probe, text)
	return args.Get(0).(orchestrationvos.Variable), args.Bool(1), args.Error(2)
}

func (m *MockVariableResolver) Interpolate(template string, variables map[string]orchestrationvos.Variable) (string, error) {
	args := m.Called(template, variables)
	return args.String(0), args.Error(1)
}

func (m *MockVariableResolver) ProcessTemplate(pathFile string, variables map[string]orchestrationvos.Variable) error {
	args := m.Called(pathFile, variables)
	return args.Error(0)
}

// Ensure MockVariableResolver implements the interface
var _ orchestrationservices.VariableResolver = (*MockVariableResolver)(nil)

type MockFingerprintService struct {
	mock.Mock
}

func (m *MockFingerprintService) CalculateCodeFingerprint() (executionstatevos.Fingerprint, error) {
	args := m.Called()
	return args.Get(0).(executionstatevos.Fingerprint), args.Error(1)
}

func (m *MockFingerprintService) CalculateStepFingerprint(stepName string) (executionstatevos.Fingerprint, error) {
	args := m.Called(stepName)
	return args.Get(0).(executionstatevos.Fingerprint), args.Error(1)
}

type MockWorkspaceManager struct {
	mock.Mock
}

func (m *MockWorkspaceManager) Prepare(stepName string) (string, error) {
	args := m.Called(stepName)
	return args.String(0), args.Error(1)
}

type MockCommandExecutor struct {
	mock.Mock
}

func (m *MockCommandExecutor) Execute(ctx context.Context, workdir, command string) (string, int, error) {
	args := m.Called(ctx, workdir, command)
	return args.String(0), args.Int(1), args.Error(2)
}

func (m *MockCommandExecutor) CreateWorkDir(paths ...string) string {
	args := m.Called(paths)
	return args.String(0)
}

type MockVarsRepository struct {
	mock.Mock
}

func (m *MockVarsRepository) FindAll() ([]orchestrationvos.Variable, error) {
	args := m.Called()
	return args.Get(0).([]orchestrationvos.Variable), args.Error(1)
}

func (m *MockVarsRepository) Save(vars []orchestrationvos.Variable) error {
	args := m.Called(vars)
	return args.Error(0)
}

type MockStateRepository struct {
	mock.Mock
}

func (m *MockStateRepository) FindStepStatus() (executionstateaggregates.StateSteps, error) {
	args := m.Called()
	return args.Get(0).(executionstateaggregates.StateSteps), args.Error(1)
}

func (m *MockStateRepository) SaveStepStatus(stateSteps executionstateaggregates.StateSteps) error {
	args := m.Called(stateSteps)
	return args.Error(0)
}

// Test Suite Setup
type OrchestrationServiceTestSuite struct {
	t *testing.T

	// Mocks
	stepVariableRepo *MockStepVariableRepository
	orderRepo        *MockOrderRepository
	scopeRepo        *MockScopeRepository
	varResolver      *MockVariableResolver
	fpService        *MockFingerprintService
	workspaceMgr     *MockWorkspaceManager
	cmdExecutor      *MockCommandExecutor
	varsRepo         *MockVarsRepository
	stateRepo        *MockStateRepository

	// Service under test
	service *OrchestrationService
}

func setup(t *testing.T) *OrchestrationServiceTestSuite {
	suite := &OrchestrationServiceTestSuite{
		t:                t,
		stepVariableRepo: new(MockStepVariableRepository),
		orderRepo:        new(MockOrderRepository),
		scopeRepo:        new(MockScopeRepository),
		varResolver:      new(MockVariableResolver),
		fpService:        new(MockFingerprintService),
		workspaceMgr:     new(MockWorkspaceManager),
		cmdExecutor:      new(MockCommandExecutor),
		varsRepo:         new(MockVarsRepository),
		stateRepo:        new(MockStateRepository),
	}

	suite.service = NewOrchestrationService(
		suite.stepVariableRepo,
		suite.orderRepo,
		suite.scopeRepo,
		suite.varResolver,
		suite.fpService,
		suite.workspaceMgr,
		suite.cmdExecutor,
		suite.varsRepo,
		suite.stateRepo,
	)
	return suite
}

func TestOrchestrationService_ExecuteOrder_Success(t *testing.T) {
	suite := setup(t)
	ctx := context.Background()
	environmentName := "staging"
	environmentValue := "stag"

	source, _ := deploymentvos.NewTemplateSource("git", "main")
	env, _ := deploymentvos.NewEnvironment(environmentName, "staging description", environmentValue)
	step1, _ := deploymententities.NewStepDefinition("test", []deploymentvos.VerificationType{deploymentvos.VerificationTypeCode}, nil)
	step2, _ := deploymententities.NewStepDefinition("supply", []deploymentvos.VerificationType{deploymentvos.VerificationTypeEnv}, nil)
	step3, _ := deploymententities.NewStepDefinition("package", []deploymentvos.VerificationType{deploymentvos.VerificationTypeEnv, deploymentvos.VerificationTypeCode}, nil)
	step4, _ := deploymententities.NewStepDefinition("deploy", []deploymentvos.VerificationType{deploymentvos.VerificationTypeEnv, deploymentvos.VerificationTypeCode}, nil)

	template, _ := deploymentaggregates.NewDeploymentTemplate(
		source,
		[]deploymentvos.Environment{env},
		[]deploymententities.StepDefinition{step1, step2, step3, step4},
	)

	product, _ := domvos.NewProduct(domvos.ProductID("productID"), "productName", "description", "team", "org")
	project, _ := domvos.NewProject(domvos.ProjectID("projectID"), "projectName", "1.0.0", "description", "team")
	domTemplate, _ := domvos.NewTemplate("url", "ref")
	technology, _ := domvos.NewTechnology("tech", "sol", "stack", "infra")
	projectDom := domaggregates.NewDeploymentObjectModel(product, project, domTemplate, technology)

	req := dto.OrderRequest{
		Ctx:            ctx,
		ProjectDom:     projectDom,
		Template:       template,
		Environment:    env,
		FinalStep:      "deploy",
		ProjectPath:    "/path/to/project",
		RepositoryName: "test-repo",
	}

	// 2. Mock Dependencies
	suite.varsRepo.On("GetStore", environmentValue).Return([]orchestrationvos.Variable{}, nil)
	history, _ := executionstateaggregates.NewScopeReceiptHistory()
	suite.scopeRepo.On("FindCodeState").Return(history, nil)
	fp, _ := executionstatevos.NewFingerprint("codefp")
	suite.fpService.On("CalculateCodeFingerprint", ctx, "/path/to/project").Return(fp, nil)
	suite.stateRepo.On("FindStateSteps", "staging").Return(executionstateaggregates.NewStateSteps(), nil)
	suite.orderRepo.On("Save", ctx, mock.AnythingOfType("*aggregates.Order"), "projectName").Return(nil)
	suite.varsRepo.On("Save", mock.Anything, "staging").Return(nil)
	suite.stateRepo.On("SaveStateSteps", mock.Anything, "staging").Return(nil)

	// Mocking for each step
	for _, step := range template.Steps() {
		stepName := step.Name()
		envHistory, _ := executionstateaggregates.NewScopeReceiptHistory()
		suite.scopeRepo.On("FindEnvironmentState", "staging", stepName).Return(envHistory, nil)
		envFp, _ := executionstatevos.NewFingerprint("envfp-" + stepName)
		suite.fpService.On("CalculateEnvironmentFingerprint", ctx, stepName, "").Return(envFp, nil)
		suite.stepVariableRepo.On("Load", "staging", stepName).Return([]orchestrationvos.Variable{}, nil)
		suite.workspaceMgr.On("PrepareStepWorkspace", "projectName", "staging", stepName, "test-repo").Return("/workdir/"+stepName, nil)
		suite.varResolver.On("Interpolate", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]vos.Variable")).Return("echo 'hello'", nil)
		suite.cmdExecutor.On("Execute", ctx, mock.AnythingOfType("string"), "echo 'hello'").Return("output", 0, nil)
		suite.cmdExecutor.On("CreateWorkDir", mock.Anything).Return("/workdir/step/cmd")

		if stepName == "test" {
			suite.scopeRepo.On("SaveCodeState", mock.AnythingOfType("*aggregates.ScopeReceiptHistory")).Return(nil)
		} else {
			suite.scopeRepo.On("SaveEnvironmentState", mock.AnythingOfType("*aggregates.ScopeReceiptHistory"), "staging", stepName).Return(nil)
		}
	}
	suite.varResolver.On("ExtractVariable", mock.Anything, mock.Anything).Return(orchestrationvos.Variable{}, false, nil)
	suite.varResolver.On("ProcessTemplate", mock.Anything, mock.Anything).Return(nil)

	// 3. Execute
	order, err := suite.service.ExecuteOrder(req)

	// 4. Assert
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, orchestrationvos.OrderStatusSuccessful, order.Status())

	// Verify that all mocks were called as expected
	suite.varsRepo.AssertExpectations(t)
	suite.scopeRepo.AssertExpectations(t)
	suite.fpService.AssertExpectations(t)
	suite.stateRepo.AssertExpectations(t)
	suite.orderRepo.AssertExpectations(t)
	suite.stepVariableRepo.AssertExpectations(t)
	suite.workspaceMgr.AssertExpectations(t)
	suite.varResolver.AssertExpectations(t)
	suite.cmdExecutor.AssertExpectations(t)
}

func TestOrchestrationService_ExecuteOrder_CommandFails(t *testing.T) {
	suite := setup(t)
	ctx := context.Background()

	// 1. Setup Request (similar to the success test, but simpler)
	source, _ := deploymentvos.NewTemplateSource("git", "main")
	env, _ := deploymentvos.NewEnvironment("staging", "staging env", "staging")
	cmd, _ := deploymentvos.NewCommandDefinition("fail-cmd", "exit 1")
	step1, _ := deploymententities.NewStepDefinition("build", nil, []deploymentvos.CommandDefinition{cmd})
	template, _ := deploymentaggregates.NewDeploymentTemplate(
		source,
		[]deploymentvos.Environment{env},
		[]deploymententities.StepDefinition{step1},
	)
	product, _ := domvos.NewProduct(domvos.ProductID("productID"), "productName", "description", "team", "org")
	project, _ := domvos.NewProject(domvos.ProjectID("projectID"), "projectName", "1.0.0", "description", "team")
	domTemplate, _ := domvos.NewTemplate("url", "ref")
	technology, _ := domvos.NewTechnology("tech", "sol", "stack", "infra")
	projectDom := domaggregates.NewDeploymentObjectModel(product, project, domTemplate, technology)

	req := dto.OrderRequest{
		Ctx:         ctx,
		ProjectDom:  projectDom,
		Template:    template,
		Environment: env,
		FinalStep:   "build",
		ProjectPath: "/path/to/project",
	}

	// 2. Mock Dependencies
	suite.varsRepo.On("GetStore", "staging").Return([]orchestrationvos.Variable{}, nil)
	history, _ := executionstateaggregates.NewScopeReceiptHistory()
	suite.scopeRepo.On("FindCodeState").Return(history, nil)
	fp, _ := executionstatevos.NewFingerprint("codefp")
	suite.fpService.On("CalculateCodeFingerprint", ctx, "/path/to/project").Return(fp, nil)
	suite.stateRepo.On("FindStateSteps", "staging").Return(executionstateaggregates.NewStateSteps(), nil)
	suite.scopeRepo.On("FindEnvironmentState", "staging", "build").Return(history, nil)
	envFp, _ := executionstatevos.NewFingerprint("envfp-build")
	suite.fpService.On("CalculateEnvironmentFingerprint", ctx, "build", "").Return(envFp, nil)
	suite.stepVariableRepo.On("Load", "staging", "build").Return([]orchestrationvos.Variable{}, nil)
	suite.workspaceMgr.On("PrepareStepWorkspace", "projectName", "staging", "build", "").Return("/workdir/build", nil)
	suite.varResolver.On("Interpolate", "exit 1", mock.AnythingOfType("map[string]vos.Variable")).Return("exit 1", nil)

	// Here's the failure injection
	suite.cmdExecutor.On("Execute", ctx, mock.AnythingOfType("string"), "exit 1").Return("error output", 1, nil)
	suite.cmdExecutor.On("CreateWorkDir", mock.Anything).Return("/workdir/step/cmd")
	suite.varResolver.On("ExtractVariable", mock.Anything, mock.Anything).Return(orchestrationvos.Variable{}, false, nil)
	suite.varResolver.On("ProcessTemplate", mock.Anything, mock.Anything).Return(nil)

	// The order should be saved even when failed
	suite.orderRepo.On("Save", ctx, mock.AnythingOfType("*aggregates.Order"), "projectName").Return(nil)

	// 3. Execute
	order, err := suite.service.ExecuteOrder(req)

	// 4. Assert
	assert.NoError(t, err) // The service itself doesn't error, it manages the failure state in the order
	assert.NotNil(t, order)
	assert.Equal(t, orchestrationvos.OrderStatusFailed, order.Status())
	assert.Equal(t, orchestrationvos.StepStatusFailed, order.StepExecutions()[0].Status())
	assert.Equal(t, orchestrationvos.CommandStatusFailed, order.StepExecutions()[0].CommandExecutions()[0].Status())

	suite.cmdExecutor.AssertExpectations(t)
	suite.orderRepo.AssertExpectations(t)
}

func TestOrchestrationService_ExecuteOrder_DependencyError(t *testing.T) {
	suite := setup(t)
	ctx := context.Background()

	// 1. Setup Request
	source, _ := deploymentvos.NewTemplateSource("git", "main")
	env, _ := deploymentvos.NewEnvironment("staging", "staging env", "staging")
	template, _ := deploymentaggregates.NewDeploymentTemplate(source, []deploymentvos.Environment{env}, nil)
	product, _ := domvos.NewProduct(domvos.ProductID("productID"), "productName", "description", "team", "org")
	project, _ := domvos.NewProject(domvos.ProjectID("projectID"), "projectName", "1.0.0", "description", "team")
	domTemplate, _ := domvos.NewTemplate("url", "ref")
	technology, _ := domvos.NewTechnology("tech", "sol", "stack", "infra")
	projectDom := domaggregates.NewDeploymentObjectModel(product, project, domTemplate, technology)

	req := dto.OrderRequest{
		Ctx:         ctx,
		ProjectDom:  projectDom,
		Template:    template,
		Environment: env,
		FinalStep:   "build",
	}

	// 2. Mock Dependencies
	suite.varsRepo.On("GetStore", "staging").Return([]orchestrationvos.Variable{}, nil)
	// Injecting dependency error
	expectedErr := fmt.Errorf("database error")
	suite.scopeRepo.On("FindCodeState").Return(nil, expectedErr)

	// 3. Execute
	order, err := suite.service.ExecuteOrder(req)

	// 4. Assert
	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.NotNil(t, order) // The order should be partially created and returned
	assert.Equal(t, orchestrationvos.OrderStatusInProgress, order.Status())

	suite.scopeRepo.AssertExpectations(t)
}

func TestOrchestrationService_ExecuteOrder_StepSkipped_Cached(t *testing.T) {
	suite := setup(t)
	ctx := context.Background()

	// 1. Setup
	source, _ := deploymentvos.NewTemplateSource("git", "main")
	env, _ := deploymentvos.NewEnvironment("staging", "staging env", "staging")
	step1, _ := deploymententities.NewStepDefinition("build", []deploymentvos.VerificationType{deploymentvos.VerificationTypeCode}, nil)
	step2, _ := deploymententities.NewStepDefinition("deploy", nil, []deploymentvos.CommandDefinition{}) // This step will run
	template, _ := deploymentaggregates.NewDeploymentTemplate(
		source,
		[]deploymentvos.Environment{env},
		[]deploymententities.StepDefinition{step1, step2},
	)
	product, _ := domvos.NewProduct(domvos.ProductID("productID"), "productName", "description", "team", "org")
	project, _ := domvos.NewProject(domvos.ProjectID("projectID"), "projectName", "1.0.0", "description", "team")
	domTemplate, _ := domvos.NewTemplate("url", "ref")
	technology, _ := domvos.NewTechnology("tech", "sol", "stack", "infra")
	projectDom := domaggregates.NewDeploymentObjectModel(product, project, domTemplate, technology)

	req := dto.OrderRequest{
		Ctx:         ctx,
		ProjectDom:  projectDom,
		Template:    template,
		Environment: env,
		FinalStep:   "deploy",
		ProjectPath: "/path/to/project",
	}

	// 2. Mocking
	// General mocks
	suite.varsRepo.On("GetStore", "staging").Return([]orchestrationvos.Variable{}, nil).Once()
	suite.orderRepo.On("Save", ctx, mock.AnythingOfType("*aggregates.Order"), "projectName").Return(nil)
	suite.varsRepo.On("Save", mock.Anything, "staging").Return(nil)
	suite.stateRepo.On("SaveStateSteps", mock.Anything, "staging").Return(nil)
	suite.varResolver.On("ExtractVariable", mock.Anything, mock.Anything).Return(orchestrationvos.Variable{}, false, nil)
	suite.varResolver.On("ProcessTemplate", mock.Anything, mock.Anything).Return(nil)

	// Mocks to make 'build' step cached
	codeFp, _ := executionstatevos.NewFingerprint("codefp")
	suite.fpService.On("CalculateCodeFingerprint", ctx, "/path/to/project").Return(codeFp, nil).Once()
	stateSteps := executionstateaggregates.NewStateSteps()
	stateStep, _ := executionstatevos.NewStateStep("build", true)
	stateSteps.AddStep(stateStep)
	suite.stateRepo.On("FindStateSteps", "staging").Return(stateSteps, nil).Once()

	codeHistory, _ := executionstateaggregates.NewScopeReceiptHistory()
	receipt, _ := executionstateaggregates.NewScopeReceipt(codeFp, executionstatevos.Fingerprint{})
	codeHistory.AddReceipt(receipt)
	suite.scopeRepo.On("FindCodeState").Return(codeHistory, nil).Once()

	// Mocks for 'deploy' step (it should run)
	suite.scopeRepo.On("FindEnvironmentState", "staging", "deploy").Return(executionstateaggregates.NewScopeReceiptHistory())
	envFp, _ := executionstatevos.NewFingerprint("envfp-deploy")
	suite.fpService.On("CalculateEnvironmentFingerprint", ctx, "deploy", "").Return(envFp, nil)
	suite.stepVariableRepo.On("Load", "staging", "deploy").Return([]orchestrationvos.Variable{}, nil)
	suite.workspaceMgr.On("PrepareStepWorkspace", "projectName", "staging", "deploy", "").Return("/workdir/deploy", nil)
	suite.varResolver.On("Interpolate", mock.AnythingOfType("string"), mock.AnythingOfType("map[string]vos.Variable")).Return("echo 'deploy'", nil)
	suite.cmdExecutor.On("Execute", ctx, mock.AnythingOfType("string"), "echo 'deploy'").Return("deployed", 0, nil)
	suite.cmdExecutor.On("CreateWorkDir", mock.Anything).Return("/workdir/step/cmd")
	suite.scopeRepo.On("SaveEnvironmentState", mock.Anything, "staging", "deploy").Return(nil)

	// 3. Execute
	order, err := suite.service.ExecuteOrder(req)

	// 4. Assert
	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, orchestrationvos.OrderStatusSuccessful, order.Status())
	assert.Equal(t, orchestrationvos.StepStatusCached, order.StepExecutions()[0].Status())     // build step
	assert.Equal(t, orchestrationvos.StepStatusSuccessful, order.StepExecutions()[1].Status()) // deploy step

	suite.cmdExecutor.AssertNumberOfCalls(t, "Execute", 1) // Only deploy step command should be executed
}

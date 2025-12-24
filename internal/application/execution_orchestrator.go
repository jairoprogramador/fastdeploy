package application

import (
	"context"
	"fmt"
	defAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/aggregates"
	defPrt "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/ports"
	defVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/vos"
	exePrt "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/ports"
	exeVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/execution/vos"
	proAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/aggregates"
	proPrt "github.com/jairoprogramador/fastdeploy-core/internal/domain/project/ports"
	staPrt "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/ports"
	staVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"
	verPrt "github.com/jairoprogramador/fastdeploy-core/internal/domain/versioning/ports"
	worAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/workspace/aggregates"
)

// ExecutionOrchestrator orquesta la ejecución de un plan completo,
// coordinando los contextos de definición, estado y ejecución.
type ExecutionOrchestrator struct {
	projectPath        string
	rootFastdeployPath string
	projectSvc         *ProjectService
	workspaceSvc       *WorkspaceService
	gitCloner          proPrt.ClonerTemplate
	versionCalculator  verPrt.VersionCalculator
	planBuilder        defPrt.PlanBuilder
	fingerprintSvc     staPrt.FingerprintService
	stateManager       staPrt.StateManager
	stepExecutor       exePrt.StepExecutor
	copyWorkdir        exePrt.CopyWorkdir
}

// NewExecutionOrchestrator crea una nueva instancia del orquestador.
func NewExecutionOrchestrator(
	projectPath string,
	rootFastdeployPath string,
	projectSvc *ProjectService,
	workspaceSvc *WorkspaceService,
	gitCloner proPrt.ClonerTemplate,
	versionCalculator verPrt.VersionCalculator,
	planBuilder defPrt.PlanBuilder,
	fingerprintSvc staPrt.FingerprintService,
	stateManager staPrt.StateManager,
	stepExecutor exePrt.StepExecutor,
	copyWorkdir exePrt.CopyWorkdir,
) *ExecutionOrchestrator {
	return &ExecutionOrchestrator{
		projectPath:        projectPath,
		rootFastdeployPath: rootFastdeployPath,
		projectSvc:         projectSvc,
		workspaceSvc:       workspaceSvc,
		gitCloner:          gitCloner,
		versionCalculator:  versionCalculator,
		planBuilder:        planBuilder,
		fingerprintSvc:     fingerprintSvc,
		stateManager:       stateManager,
		stepExecutor:       stepExecutor,
		copyWorkdir:        copyWorkdir,
	}
}

// ExecutePlan es el caso de uso principal que ejecuta un plan de despliegue.
func (o *ExecutionOrchestrator) ExecutePlan(ctx context.Context, stepName, envName string) error {
	// 1. Inicializar, Cargar y Clonar
	project, err := o.loadProject(ctx, o.projectPath)
	if err != nil {
		return err
	}
	workspace, err := o.loadWorkspace(project, o.rootFastdeployPath)
	if err != nil {
		return err
	}

	templateLocalPath := workspace.TemplatePath()
	err = o.cloneTemplate(ctx, project, templateLocalPath)
	if err != nil {
		return err
	}

	planDef, err := o.buildPlan(ctx, templateLocalPath, stepName, envName)
	if err != nil {
		return err
	}

	version, commit, err := o.versionCalculator.CalculateNextVersion(ctx, o.projectPath, false)
	if err != nil {
		return err
	}

	environment := planDef.Environment().String()

	projectVars := o.prepareProjectVariables(project)
	othersVars := o.prepareOthersVariables(
		environment, o.projectPath, version.String(), commit.String())

	cumulativeVars := make(exeVos.VariableSet)
	cumulativeVars.AddAll(projectVars)
	cumulativeVars.AddAll(othersVars)

	fmt.Println("Iniciando la ejecución del plan...")

	// 3. Bucle de Ejecución Paso a Paso
	for _, stepDef := range planDef.Steps() {

		// 3a. Comprobación de Estado (Cache Check)
		fingerprints, err := o.generateStepFingerprints(o.projectPath, environment, workspace, stepDef.NameDef())
		if err != nil {
			return fmt.Errorf("error al generar fingerprint para el paso '%s': %w", stepDef.NameDef().Name(), err)
		}

		stateTablePath, err := workspace.StateTablePath(stepDef.NameDef().Name())
		if err != nil {
			return fmt.Errorf("error al obtener la ruta del estado del paso '%s': %w", stepDef.NameDef().Name(), err)
		}
		hasChanged, err := o.stateManager.HasStateChanged(stateTablePath, fingerprints, staVos.NewCachePolicy(0))
		if err != nil {
			return fmt.Errorf("error al comprobar el estado del paso '%s': %w", stepDef.NameDef().Name(), err)
		}

		if !hasChanged {
			fmt.Printf("  - Paso '%s' no ha cambiado (cache HIT). Omitiendo.\n", stepDef.NameDef().Name())
			continue // Saltar al siguiente paso
		}

		stepPathEnv := workspace.ScopeWorkdirPath(planDef.Environment().String(), stepDef.NameDef().Name())
		err = o.copyWorkdir.Copy(ctx, workspace.StepTemplatePath(stepDef.NameDef().FullName()), stepPathEnv)
		if err != nil {
			return fmt.Errorf("error al copiar el paso '%s' al workspace: %w", stepPathEnv, err)
		}

		execStep, err := mapToExecutionStep(stepDef, stepPathEnv)
		if err != nil {
			return fmt.Errorf("error al mapear la definición del paso '%s': %w", stepDef.NameDef().Name(), err)
		}

		// 3b. Ejecución del Paso
		fmt.Printf("  - Paso '%s' ha cambiado. Ejecutando...\n", stepDef.NameDef().Name())
		execResult, err := o.stepExecutor.Execute(ctx, execStep, cumulativeVars)
		if err != nil {
			return fmt.Errorf("la ejecución del paso '%s' falló: %w", stepDef.NameDef().Name(), err)
		}
		if execResult.Error != nil || execResult.Status == exeVos.Failure {
			fmt.Println("--- Logs del fallo ---")
			fmt.Println(execResult.Logs)
			fmt.Println("--------------------")
			return fmt.Errorf("el paso '%s' finalizó con error: %w", stepDef.NameDef().Name(), execResult.Error)
		}

		// 3c. Actualización de Variables y Estado
		fmt.Printf("Paso '%s' completado. Salida:\n%s\n", stepDef.NameDef().Name(), execResult.Logs)
		for key, value := range execResult.OutputVars {
			cumulativeVars[key] = value
		}

		if err := o.stateManager.UpdateState(stateTablePath, fingerprints); err != nil {
			// Esto es una advertencia. El flujo principal fue exitoso, pero el estado no se guardó.
			fmt.Printf("ADVERTENCIA: no se pudo guardar el estado del paso '%s'. Se re-ejecutará la próxima vez. Error: %v\n", stepDef.NameDef().Name(), err)
		}
	}

	fmt.Println("\n¡Ejecución completada con éxito!")
	return nil
}

func (o *ExecutionOrchestrator) loadProject(ctx context.Context, projectPath string) (*proAgg.Project, error) {
	// 1. Cargar el Proyecto
	project, err := o.projectSvc.Load(ctx, projectPath)
	if err != nil {
		return nil, fmt.Errorf("error al cargar el proyecto: %w", err)
	}
	return project, nil
}

func (o *ExecutionOrchestrator) loadWorkspace(project *proAgg.Project, rootFastdeployPath string) (*worAgg.Workspace, error) {
	// 2. Crear el Workspace
	workspace, err := o.workspaceSvc.NewWorkspace(
		rootFastdeployPath, project.Data().Name(), project.TemplateRepo().DirName())
	if err != nil {
		return nil, fmt.Errorf("error al cargar el workspace: %w", err)
	}
	return workspace, nil
}

func (o *ExecutionOrchestrator) cloneTemplate(
	ctx context.Context, project *proAgg.Project, templateLocalPath string) error {
	// 3. Asegurar que el template está clonado
	err := o.gitCloner.EnsureCloned(ctx, project.TemplateRepo().URL(),
		project.TemplateRepo().Ref(), templateLocalPath)
	if err != nil {
		return fmt.Errorf("no se pudo clonar el repositorio de plantillas: %w", err)
	}
	return nil
}

func (o *ExecutionOrchestrator) buildPlan(
	ctx context.Context, templateLocalPath, stepName, envName string) (*defAgg.ExecutionPlanDefinition, error) {

	// 4. Cargar la definición del plan desde el template clonado
	planDef, err := o.planBuilder.Build(ctx, templateLocalPath, stepName, envName)
	if err != nil {
		return nil, fmt.Errorf("error al cargar la definición: %w", err)
	}

	return planDef, nil
}

func (o *ExecutionOrchestrator) prepareProjectVariables(project *proAgg.Project) exeVos.VariableSet {
	vars := make(exeVos.VariableSet)
	vars.Add("project_id", project.ID().String()[:8])
	vars.Add("project_name", project.Data().Name())
	vars.Add("project_organization", project.Data().Organization())
	vars.Add("project_team", project.Data().Team())
	//vars.Add("project_version", project.Data().Version())
	return vars
}

func (o *ExecutionOrchestrator) prepareOthersVariables(environment, projectWorkdir, version, commit string) exeVos.VariableSet {
	vars := make(exeVos.VariableSet)
	vars.Add("project_version", version)
	vars.Add("project_revision", commit)
	vars.Add("environment", environment)
	vars.Add("project_workdir", projectWorkdir)
	vars.Add("tool_name", "fastdeploy")
	return vars
}

func (o *ExecutionOrchestrator) generateCodeFingerprint(projectPath string) (staVos.Fingerprint, error) {
	codeFp, err := o.fingerprintSvc.FromDirectory(projectPath)
	if err != nil {
		return staVos.Fingerprint{}, fmt.Errorf("no se pudo generar el fingerprint para el proyecto: %w", err)
	}
	return codeFp, nil
}

func (o *ExecutionOrchestrator) generateInstructionFingerprint(templateInstPath string) (staVos.Fingerprint, error) {
	codeFp, err := o.fingerprintSvc.FromDirectory(templateInstPath)
	if err != nil {
		return staVos.Fingerprint{}, fmt.Errorf("no se pudo generar el fingerprint para las instrucciones: %w", err)
	}
	return codeFp, nil
}

func (o *ExecutionOrchestrator) generateVarsFingerprint(templateVarsPath string) (staVos.Fingerprint, error) {
	codeFp, err := o.fingerprintSvc.FromFile(templateVarsPath)
	if err != nil {
		return staVos.Fingerprint{}, fmt.Errorf("no se pudo generar el fingerprint para las variables: %w", err)
	}
	return codeFp, nil
}

func (o *ExecutionOrchestrator) generateStepFingerprints(
	projectPath, environment string,
	workspace *worAgg.Workspace,
	stepDef defVos.StepNameDefinition) (staVos.CurrentStateFingerprints, error) {

	envFp, err := staVos.NewEnvironment(environment)
	if err != nil {
		return staVos.CurrentStateFingerprints{}, err
	}

	codeFp, err := o.generateCodeFingerprint(projectPath)
	if err != nil {
		return staVos.CurrentStateFingerprints{}, err
	}

	instructionPath := workspace.StepTemplatePath(stepDef.FullName())
	instFp, err := o.generateInstructionFingerprint(instructionPath)
	if err != nil {
		return staVos.CurrentStateFingerprints{}, err
	}

	varsPath := workspace.VarsTemplatePath(stepDef.Name(), environment)
	varsFp, err := o.generateVarsFingerprint(varsPath)
	if err != nil {
		return staVos.CurrentStateFingerprints{}, err
	}

	return staVos.NewCurrentStateFingerprints(codeFp, instFp, varsFp, envFp), nil
}

package application

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/newinternal/application/ports"
	deployment "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/aggregates"
	deploymentvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	deploymententities "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/entities"
	deploymentports "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/ports"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/aggregates"
	orchestrationports "github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/ports"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/services"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/application/dto"
)

// OrchestrationService es el servicio de aplicación que coordina el caso de uso
// principal: la ejecución de una orden de despliegue.
type OrchestrationService struct {
	templateRepo deploymentports.TemplateRepository
	orderRepo    orchestrationports.OrderRepository
	stepVarRepo  orchestrationports.StepVariableRepository
	varResolver  services.VariableResolver
	cmdExecutor  ports.CommandExecutor
}

// NewOrchestrationService crea una nueva instancia de OrchestrationService.
// Las dependencias se inyectan siguiendo el principio de Inversión de Dependencias.
func NewOrchestrationService(
	templateRepo deploymentports.TemplateRepository,
	orderRepo orchestrationports.OrderRepository,
	stepVarRepo orchestrationports.StepVariableRepository,
	varResolver services.VariableResolver,
	cmdExecutor ports.CommandExecutor,
) *OrchestrationService {
	return &OrchestrationService{
		templateRepo: templateRepo,
		orderRepo:    orderRepo,
		stepVarRepo:  stepVarRepo,
		varResolver:  varResolver,
		cmdExecutor:  cmdExecutor,
	}
}

// ExecuteOrder orquesta todo el proceso de ejecución de un despliegue.
func (s *OrchestrationService) ExecuteOrder(req dto.ExecuteOrderRequest) (*aggregates.Order, error) {
	// 1. Obtener la "receta" del despliegue.
	template, err := s.templateRepo.GetTemplate(req.Ctx, req.TemplateSource)
	if err != nil {
		return nil, fmt.Errorf("error al obtener la plantilla de despliegue: %w", err)
	}

	// 2. Validar y obtener el ambiente objetivo.
	targetEnv, err := findEnvironment(template.Environments(), req.EnvironmentName)
	if err != nil {
		return nil, err
	}

	// 3. Determinar los pasos a ejecutar y cargar todas sus variables.
	stepsToExecute := getStepsToExecute(template, req.FinalStepName)
	allVars := req.InitialVariables
	for _, stepDef := range stepsToExecute {
		stepVars, err := s.stepVarRepo.Load(req.Ctx, targetEnv, stepDef)
		if err != nil {
			return nil, fmt.Errorf("error al cargar las variables para el paso '%s': %w", stepDef.Name(), err)
		}
		allVars = append(allVars, stepVars...)
	}

	// TODO: Procesar variables que dependen de otras.

	// 4. Crear el agregado Order. El dominio se encarga de la lógica de negocio.
	order, err := aggregates.NewOrder(
		vos.NewOrderID(),
		template,
		targetEnv,
		req.FinalStepName,
		req.SkippedStepNames,
		allVars,
	)
	if err != nil {
		return nil, fmt.Errorf("error al crear la orden de ejecución: %w", err)
	}

	// 5. Bucle de ejecución de la orden (pasos y comandos).
	// Esta lógica se implementaría aquí, iterando sobre los pasos de la orden,
	// ejecutando comandos con s.cmdExecutor, y actualizando el agregado Order
	// después de cada paso con s.orderRepo.Save(order).

	// Por ahora, devolvemos la orden recién creada.
	return order, nil
}

// findEnvironment es una función de ayuda para buscar un ambiente por su nombre.
func findEnvironment(envs []deploymentvos.Environment, name string) (deploymentvos.Environment, error) {
	for _, env := range envs {
		if env.Name() == name || env.Value() == name {
			return env, nil
		}
	}
	return deploymentvos.Environment{}, fmt.Errorf("el ambiente '%s' no se encontró en la definición de la plantilla", name)
}

// getStepsToExecute es una función de ayuda para determinar la secuencia de pasos.
func getStepsToExecute(template *deployment.DeploymentTemplate, finalStepName string) []deploymententities.StepDefinition {
	// Lógica para encontrar el índice del paso final y devolver todos los pasos hasta ese punto.
	return []deploymententities.StepDefinition{} // Implementación pendiente.
}

package state

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/aggregates"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/entities"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/state/dto"
)

// FileOrderRepository implementa la interfaz ports.OrderRepository utilizando el sistema de archivos.
// Persiste el estado de cada orden como un archivo YAML separado.
type FileOrderRepository struct {
	basePath string
}

// NewFileOrderRepository crea una nueva instancia del repositorio de órdenes.
func NewFileOrderRepository(basePath string) (*FileOrderRepository, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("no se pudo crear el directorio base para el estado de las órdenes: %w", err)
	}
	return &FileOrderRepository{basePath: basePath}, nil
}

// Save serializa el agregado Order a un archivo YAML.
func (r *FileOrderRepository) Save(_ context.Context, order *aggregates.Order) error {
	orderDTO := r.mapDomainToDTO(order)

	data, err := yaml.Marshal(orderDTO)
	if err != nil {
		return fmt.Errorf("error al serializar la orden a YAML: %w", err)
	}

	filePath := filepath.Join(r.basePath, fmt.Sprintf("%s.yaml", order.ID().String()))
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("error al guardar el archivo de estado de la orden: %w", err)
	}

	return nil
}

// FindByID lee un archivo YAML y lo deserializa para reconstruir un agregado Order.
func (r *FileOrderRepository) FindByID(_ context.Context, id vos.OrderID) (*aggregates.Order, error) {
	filePath := filepath.Join(r.basePath, fmt.Sprintf("%s.yaml", id.String()))
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no se encontró la orden con ID %s: %w", id.String(), err)
		}
		return nil, fmt.Errorf("error al leer el archivo de estado de la orden: %w", err)
	}

	var orderDTO dto.OrderDTO
	if err := yaml.Unmarshal(data, &orderDTO); err != nil {
		return nil, fmt.Errorf("error al deserializar el estado de la orden desde YAML: %w", err)
	}

	return r.mapDTOToDomain(orderDTO)
}

// mapDomainToDTO convierte el agregado de dominio a su DTO.
func (r *FileOrderRepository) mapDomainToDTO(order *aggregates.Order) dto.OrderDTO {
	stepDTOs := make([]dto.StepExecutionDTO, 0, len(order.StepExecutions()))
	for _, step := range order.StepExecutions() {
		cmdDTOs := make([]dto.CommandExecutionDTO, 0, len(step.CommandExecutions()))
		for _, cmd := range step.CommandExecutions() {
			cmdDTOs = append(cmdDTOs, dto.CommandExecutionDTO{
				Name:         cmd.Name(),
				Definition:   cmd.Definition(), // <-- Guardamos la definición
				Status:       cmd.Status(),
				ResolvedCmd:  cmd.ResolvedCmd(),
				ExecutionLog: cmd.ExecutionLog(),
				OutputVars:   cmd.OutputVars(),
			})
		}
		stepDTOs = append(stepDTOs, dto.StepExecutionDTO{
			Name:              step.Name(),
			Status:            step.Status(),
			CommandExecutions: cmdDTOs,
		})
	}

	return dto.OrderDTO{
		ID:                order.ID().String(),
		Status:            order.Status(),
		TargetEnvironment: order.TargetEnvironment(),
		StepExecutions:    stepDTOs,
		VariableMap:       order.VariableMap(),
	}
}

// mapDTOToDomain reconstruye el agregado de dominio desde su DTO.
func (r *FileOrderRepository) mapDTOToDomain(dto dto.OrderDTO) (*aggregates.Order, error) {
	orderID, err := vos.OrderIDFromString(dto.ID)
	if err != nil {
		return nil, err
	}

	stepExecs := make([]*entities.StepExecution, 0, len(dto.StepExecutions))
	for _, stepDTO := range dto.StepExecutions {
		cmdExecs := make([]*entities.CommandExecution, 0, len(stepDTO.CommandExecutions))
		for _, cmdDTO := range stepDTO.CommandExecutions {
			cmdExec := entities.RehydrateCommandExecution(
				cmdDTO.Name,
				cmdDTO.Status,
				cmdDTO.Definition, // <-- Rehidratamos usando la definición del DTO
				cmdDTO.ResolvedCmd,
				cmdDTO.ExecutionLog,
				cmdDTO.OutputVars,
			)
			cmdExecs = append(cmdExecs, cmdExec)
		}
		stepExecs = append(stepExecs, entities.RehydrateStepExecution(
			stepDTO.Name, stepDTO.Status, cmdExecs,
		))
	}

	order := aggregates.RehydrateOrder(
		orderID,
		dto.Status,
		dto.TargetEnvironment,
		stepExecs,
		dto.VariableMap,
	)
	return order, nil
}

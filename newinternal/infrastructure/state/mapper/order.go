package mapper

import (
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/aggregates"
	//"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/entities"
	//"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
	//"github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/infrastructure/state/dto"
)

/* func OrderToDomain(dto dto.OrderDTO) (*aggregates.Order, error) {
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
				vos.CommandStatusFromString(cmdDTO.Status),
				cmdDTO.Definition, // <-- Rehidratamos usando la definición del DTO
				cmdDTO.ResolvedCmd,
				cmdDTO.ExecutionLog,
				cmdDTO.OutputVars,
			)
			cmdExecs = append(cmdExecs, cmdExec)
		}
		stepExecs = append(stepExecs, entities.RehydrateStepExecution(
			stepDTO.Name, vos.StepStatusFromString(stepDTO.Status), cmdExecs,
		))
	}

	order := aggregates.RehydrateOrder(
		orderID,
		vos.OrderStatusFromString(dto.Status),
		deploymentvos.RehydrateEnvironment(dto.TargetEnvironment),
		stepExecs,
		dto.VariableMap,
	)
	return order, nil
} */

func OrderToDTO(order *aggregates.Order) dto.OrderDTO {
	stepDTOs := make([]dto.StepExecutionDTO, 0, len(order.StepExecutions()))
	for _, step := range order.StepExecutions() {
		cmdDTOs := make([]dto.CommandExecutionDTO, 0, len(step.CommandExecutions()))
		for _, cmd := range step.CommandExecutions() {

			outputVars := make(map[string]string)
			for _, output := range cmd.OutputVars() {
				outputVars[output.Key()] = output.Value()
			}

			cmdDTOs = append(cmdDTOs, dto.CommandExecutionDTO{
				Name:         cmd.Name(),
				//Definition:   cmd.Definition(), // <-- Guardamos la definición
				Status:       cmd.Status().String(),
				ResolvedCmd:  cmd.ResolvedCmd(),
				ExecutionLog: cmd.ExecutionLog(),
				OutputVars:   outputVars,
			})
		}
		stepDTOs = append(stepDTOs, dto.StepExecutionDTO{
			Name:              step.Name(),
			Status:            step.Status().String(),
			CommandExecutions: cmdDTOs,
		})
	}

	variableMap := make(map[string]string)
	for _, variable := range order.VariableMap() {
		variableMap[variable.Key()] = variable.Value()
	}

	return dto.OrderDTO{
		ID:                order.ID().String(),
		Status:            order.Status().String(),
		TargetEnvironment: order.TargetEnvironment().Name(),
		StepExecutions:    stepDTOs,
		VariableMap:       variableMap,
	}
}

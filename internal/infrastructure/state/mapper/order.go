package mapper

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/aggregates"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/state/dto"
)

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

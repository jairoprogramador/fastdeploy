package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/orchestration/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/orchestration/dto"
)

func OrderToDTO(order *aggregates.Order) dto.OrderDTO {
	stepDTOs := make([]dto.StepRecordDTO, 0, len(order.StepsRecord()))
	for _, step := range order.StepsRecord() {
		stepDTOs = append(stepDTOs, StepToDTO(step))
	}

	variableMap := make(map[string]string)
	for _, variable := range order.Outputs() {
		variableMap[variable.Name()] = variable.Value()
	}

	return dto.OrderDTO{
		ID:                order.ID().String(),
		Status:            order.Status().String(),
		TargetEnvironment: order.Environment(),
		StepRecords:       stepDTOs,
		VariableMap:       variableMap,
	}
}

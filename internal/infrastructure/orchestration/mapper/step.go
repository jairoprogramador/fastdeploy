package mapper

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/entities"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/orchestration/dto"
)

func StepToDTO(step *entities.StepRecord) dto.StepRecordDTO {
	cmdDTOs := make([]dto.CommandRecordDTO, 0, len(step.Commands()))

	for _, cmd := range step.Commands() {
		cmdDTOs = append(cmdDTOs, CommandToDTO(cmd))
	}

	return dto.StepRecordDTO{
		Name:              step.Name(),
		Status:            step.Status().String(),
		CommandExecutions: cmdDTOs,
	}
}

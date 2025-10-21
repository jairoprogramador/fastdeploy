package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/orchestration/entities"
	"github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/orchestration/dto"
)

func CommandToDTO(cmd *entities.CommandRecord) dto.CommandRecordDTO {
	outputVars := make(map[string]string)
	for _, output := range cmd.Outputs() {
		outputVars[output.Name()] = output.Value()
	}

	return dto.CommandRecordDTO{
		Name:        cmd.Name(),
		Status:      cmd.Status().String(),
		ResolvedCmd: cmd.Command(),
		Record:      cmd.Record(),
		OutputVars:  outputVars,
	}
}

package mapper

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy/internal/domain/orchestration/vos"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/git/dto"
)

// VariablesToDomain convierte un slice de DTOs de variables a un slice de VOs del dominio.
func VariablesToDomain(dtos []dto.VariableDTO) ([]vos.Variable, error) {
	variables := make([]vos.Variable, 0, len(dtos))
	for _, dto := range dtos {
		variable, err := vos.NewVariable(dto.Name, dto.Value)
		if err != nil {
			return nil, fmt.Errorf("variable inválida en el archivo de configuración (name: %s): %w", dto.Name, err)
		}
		variables = append(variables, variable)
	}
	return variables, nil
}

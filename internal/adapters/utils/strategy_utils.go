package utils

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/executor"
	"github.com/jairoprogramador/fastdeploy/internal/constants"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/context"
	"path/filepath"
)

func ExecuteStepFromFile(
	ctx context.Context,
	repositoryPath string,
	stepName string,
	executor executor.ExecutorCmd,
) error {
	technology, err := ctx.Get(constants.Technology)
	if err != nil {
		return fmt.Errorf("no se pudo obtener la tecnología del proyecto: %w", err)
	}

	repositoryFilePath := filepath.Join(repositoryPath, stepName, technology, constants.CommandFileName)

	if err := executor.Execute(repositoryFilePath); err != nil {
		return fmt.Errorf("error al ejecutar el paso '%s': %w", stepName, err)
	}

	return nil
}

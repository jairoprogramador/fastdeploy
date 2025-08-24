package strategy

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/jairoprogramador/fastdeploy/internal/domain/deployment"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
)

type BaseStrategy struct {
	Executor service.ExecutorCmd
}

func (s *BaseStrategy) ExecuteStep(
	ctx deployment.Context,
	stepName string,
	executor service.ExecutorCmd,
) error {
	homeDirPath, err := s.getHomeDirPath()
	if err != nil {
		return err
	}

	nameRepository, _ := ctx.Get(constants.KeyNameRepository)
	technologyName, _ := ctx.Get(constants.KeyTechnologyName)
	technologyVersion, _ := ctx.Get(constants.KeyTechnologyVersion)

	repositoryFilePath := filepath.Join(homeDirPath, nameRepository, 
		stepName, technologyName, technologyVersion, constants.CommandFileName)

	if err := executor.Execute(repositoryFilePath); err != nil {
		return err
	}

	return nil
}

func (s *BaseStrategy) getHomeDirPath() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio del usuario: %w", err)
	}

	return filepath.Join(currentUser.HomeDir, constants.FastDeployDir), nil
}

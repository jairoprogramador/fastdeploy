package strategy

/* import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	values "github.com/jairoprogramador/fastdeploy/internal/domain/context/values"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/constants"
	"github.com/jairoprogramador/fastdeploy/internal/infrastructure/deployment/executor/service"
)

type BaseStrategy struct {
	Executor service.ExecutorCmd
}

func (s *BaseStrategy) ExecuteStep(
	ctx *values.ContextValue,
	stepName string,
	executor service.ExecutorCmd,
) error {
	homeDirPath, err := s.getHomeDirPath()
	if err != nil {
		return err
	}

	technologyName, _ := ctx.Get(constants.ProjectTechnology)
	nameRepository, _ := ctx.Get(constants.DeploymentRepositoryName)

	var repositoryFilePath string
	if technologyName == "" {
		repositoryFilePath = filepath.Join(homeDirPath, nameRepository,
			constants.RepositoryStepsDir, stepName, constants.CommandFileName)
	} else {
		repositoryFilePath = filepath.Join(homeDirPath, nameRepository,
			constants.RepositoryStepsDir, technologyName, stepName, constants.CommandFileName)
	}

	if err := executor.Execute(repositoryFilePath, ctx); err != nil {
		return err
	}

	return nil
}

func (s *BaseStrategy) getHomeDirPath() (string, error) {
	if fastDeployHome := os.Getenv("FASTDEPLOY_HOME"); fastDeployHome != "" {
		return fastDeployHome, nil
	}

	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el directorio del usuario: %w", err)
	}
	return filepath.Join(currentUser.HomeDir, constants.FastDeployDir), nil
}
 */
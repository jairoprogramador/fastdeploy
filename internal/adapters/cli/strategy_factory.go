package cli

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/adapters/cli/strategies"
	domain "github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies"
)

func GetStrategyFactory(projectTechnology string, repositoryPath string) (domain.StrategyFactory, error) {
	switch projectTechnology {
	case "java":
		return strategies.NewJavaFactory(repositoryPath), nil
	case "node":
		return strategies.NewNodeFactory(repositoryPath), nil
	default:
		return nil, fmt.Errorf("tecnología de proyecto no soportada: %s", projectTechnology)
	}
}

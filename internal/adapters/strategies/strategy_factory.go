package strategies

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies"
)

func GetStrategyFactory(projectTechnology string, repositoryPath string) (strategies.StrategyFactory, error) {
	switch projectTechnology {
	case "java":
		return NewJavaFactory(repositoryPath), nil
	case "node":
		return NewNodeFactory(repositoryPath), nil
	default:
		return nil, fmt.Errorf("tecnología de proyecto no soportada: %s", projectTechnology)
	}
}

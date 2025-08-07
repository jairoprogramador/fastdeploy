package cli

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/core/domain/strategies"
)

func GetStrategyFactory(projectTechnology string) (strategies.StrategyFactory, error) {
	switch projectTechnology {
	case "java":
		return &strategies.JavaFactory{}, nil
	case "node":
		return &strategies.NodeFactory{}, nil
	default:
		return nil, fmt.Errorf("tecnología de proyecto no soportada: %s", projectTechnology)
	}
}
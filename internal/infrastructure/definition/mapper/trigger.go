package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/definition/vos"
)

func TriggersToDomain(scopes []string) []vos.TriggerDefinition {
	triggers := make([]vos.TriggerDefinition, 0, len(scopes))
	for _, scope := range scopes {
		triggers = append(triggers, vos.TriggerFromString(scope))
	}
	return triggers
}

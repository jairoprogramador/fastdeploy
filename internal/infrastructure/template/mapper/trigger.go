package mapper

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/template/vos"
)

func TriggersToDomain(scopes []string) []vos.Trigger {
	triggers := make([]vos.Trigger, 0, len(scopes))
	for _, scope := range scopes {
		triggers = append(triggers, vos.TriggerFromString(scope))
	}
	return triggers
}

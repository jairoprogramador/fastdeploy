package services

import (
	"fmt"
	"github.com/jairoprogramador/fastdeploy/internal/domain/state/aggregates"
	"github.com/jairoprogramador/fastdeploy/internal/domain/state/vos"
)

type ExecutionPolicyService struct {
}

func NewExecutionPolicyService() ExecutionPolicyService {
	return ExecutionPolicyService{}
}

func (s ExecutionPolicyService) Decide(
	lastState *aggregates.ExecutionState,
	triggers []int,
	currentFingerprints map[vos.Trigger]vos.Fingerprint,
) vos.ExecutionDecision {

	if lastState == nil {
		return vos.Execute("Primera ejecución del paso.")
	}

	if len(triggers) == 0 {
		return vos.Execute("El paso no tiene triggers configurados.")
	}

	for _, trigger := range triggers {
		lastFingerprint, okLast := lastState.GetFingerprint(vos.NewTrigger(trigger))
		currentFingerprint, okCurrent := currentFingerprints[vos.NewTrigger(trigger)]

		if !okLast || !okCurrent {
			return vos.Execute(fmt.Sprintf("El trigger '%d' no tiene fingerprint configurado.", trigger))
		}
		if !lastFingerprint.Equals(currentFingerprint) {
			return vos.Execute(fmt.Sprintf("Cambio detectado en el trigger '%d'.", trigger))
		}
	}
	return vos.Skip("No se han detectado cambios en los triggers.")
}
package services

import (
	"fmt"

	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/aggregates"
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"
)

type FingerprintPolicyService struct {
}

func NewFingerprintPolicyService() FingerprintPolicyService {
	return FingerprintPolicyService{}
}

func (s FingerprintPolicyService) Decide(
	lastFingerprints *aggregates.FingerprintState,
	triggers []int,
	currentFingerprints *aggregates.FingerprintState,
) vos.FingerprintDecision {

	if lastFingerprints == nil {
		return vos.Execute("Primera ejecución del paso")
	}

	if len(triggers) == 0 {
		return vos.Execute("El paso no tiene triggers configurados")
	}

	for _, trigger := range triggers {
		triggerVO := vos.NewTrigger(trigger)

		lastFingerprint, okLast := lastFingerprints.GetFingerprint(triggerVO)
		currentFingerprint, okCurrent := currentFingerprints.GetFingerprint(triggerVO)

		if !okLast || !okCurrent {
			return vos.Execute(fmt.Sprintf("El trigger '%s' no tiene fingerprint configurado", triggerVO.String()))
		}
		if !lastFingerprint.Equals(currentFingerprint) {
			return vos.Execute(fmt.Sprintf("Cambio detectado en el trigger '%s'", triggerVO.String()))
		}
	}
	return vos.Skip("No se han detectado cambios en los triggers")
}

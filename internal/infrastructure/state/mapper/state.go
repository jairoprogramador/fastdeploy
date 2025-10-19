package mapper

import (
	staVos "github.com/jairoprogramador/fastdeploy/internal/domain/state/vos"
	staAgg "github.com/jairoprogramador/fastdeploy/internal/domain/state/aggregates"

	staDto "github.com/jairoprogramador/fastdeploy/internal/infrastructure/state/dto"
)

func ToDTO(state *staAgg.ExecutionState) staDto.StateFingerprintDTO {
	fingerprints := make(map[int]string)
	for trigger, fingerprint := range state.Fingerprints() {
		fingerprints[trigger.Int()] = fingerprint.String()
	}
	return staDto.StateFingerprintDTO{
		StepName: state.StepName(),
		Fingerprints: fingerprints,
	}
}

func ToDomain(dto staDto.StateFingerprintDTO) *staAgg.ExecutionState {
	executionState := staAgg.NewExecutionState(dto.StepName)

	for trigger, fingerprint := range dto.Fingerprints {
		if trigger >= 0 && fingerprint != "" {
			fingerprint, _ := staVos.NewFingerprint(fingerprint)
			triggerVO := staVos.NewTrigger(trigger)
			executionState.SetFingerprint(triggerVO, fingerprint)
		}
	}

	return executionState
}
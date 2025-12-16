package mapper

import (
	staAgg "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/aggregates"
	staVos "github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"

	staDto "github.com/jairoprogramador/fastdeploy-core/internal/infrastructure/state/dto"
)

func ToDTO(state *staAgg.FingerprintState) staDto.FingerprintStateDTO {
	fingerprints := make(map[int]string)
	for trigger, fingerprint := range state.Fingerprints() {
		fingerprints[trigger.Int()] = fingerprint.String()
	}
	return staDto.FingerprintStateDTO{
		StepName:     state.StepName(),
		Fingerprints: fingerprints,
	}
}

func ToDomain(dto staDto.FingerprintStateDTO) *staAgg.FingerprintState {
	fingerprints := make(map[staVos.Trigger]staVos.Fingerprint)
	for trigger, fingerprint := range dto.Fingerprints {
		if trigger >= 0 && fingerprint != "" {
			newFingerprint, _ := staVos.NewFingerprint(fingerprint)
			newTrigger := staVos.NewTrigger(trigger)
			fingerprints[newTrigger] = newFingerprint
		}
	}
	return staAgg.HydrateFingerprintState(dto.StepName, fingerprints)
}

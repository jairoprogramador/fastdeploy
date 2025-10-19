package aggregates

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/state/vos"
)

type FingerprintState struct {
	stepName     string
	fingerprints map[vos.Trigger]vos.Fingerprint
}

func NewFingerprintState(stepName string) *FingerprintState {
	return &FingerprintState{
		stepName:     stepName,
		fingerprints: make(map[vos.Trigger]vos.Fingerprint),
	}
}

func (s *FingerprintState) SetFingerprint(trigger vos.Trigger, fingerprint vos.Fingerprint) {
	s.fingerprints[trigger] = fingerprint
}

func (s *FingerprintState) GetFingerprint(trigger vos.Trigger) (vos.Fingerprint, bool) {
	fingerprint, ok := s.fingerprints[trigger]
	return fingerprint, ok
}

func (s *FingerprintState) ExistsFingerprint(trigger vos.Trigger) bool {
	fingerprint, ok := s.GetFingerprint(trigger)
	return ok && !fingerprint.IsZero()
}

func (s *FingerprintState) StepName() string {
	return s.stepName
}

func (s *FingerprintState) Fingerprints() map[vos.Trigger]vos.Fingerprint {
	return s.fingerprints
}

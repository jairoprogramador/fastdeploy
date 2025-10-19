package aggregates

import (
	"github.com/jairoprogramador/fastdeploy/internal/domain/state/vos"
)

type ExecutionState struct {
	stepName   string
	fingerprints map[vos.Trigger]vos.Fingerprint
}

func NewExecutionState(stepName string) *ExecutionState {
	return &ExecutionState{
		stepName:   stepName,
		fingerprints: make(map[vos.Trigger]vos.Fingerprint),
	}
}

func (s *ExecutionState) SetFingerprint(trigger vos.Trigger, fingerprint vos.Fingerprint) {
	s.fingerprints[trigger] = fingerprint
}

func (s *ExecutionState) GetFingerprint(trigger vos.Trigger) (vos.Fingerprint, bool) {
	fingerprint, ok := s.fingerprints[trigger]
	return fingerprint, ok
}

func (s *ExecutionState) ExistsFingerprint(trigger vos.Trigger) bool {
	fingerprint, ok := s.GetFingerprint(trigger)
	return ok && !fingerprint.IsZero()
}

func (s *ExecutionState) StepName() string {
	return s.stepName
}

func (s *ExecutionState) Fingerprints() map[vos.Trigger]vos.Fingerprint {
	return s.fingerprints
}
package aggregates

import (
	"github.com/jairoprogramador/fastdeploy-core/internal/domain/state/vos"
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

func HydrateFingerprintState(stepName string, fingerprints map[vos.Trigger]vos.Fingerprint) *FingerprintState {
	return &FingerprintState{
		stepName:     stepName,
		fingerprints: fingerprints,
	}
}

func (s *FingerprintState) AddFingerprintdCode(fingerprint vos.Fingerprint) {
	s.fingerprints[vos.ScopeCode] = fingerprint
}

func (s *FingerprintState) AddFingerprintdRecipe(fingerprint vos.Fingerprint) {
	s.fingerprints[vos.ScopeRecipe] = fingerprint
}

func (s *FingerprintState) AddFingerprintdVars(fingerprint vos.Fingerprint) {
	s.fingerprints[vos.ScopeVars] = fingerprint
}

func (s *FingerprintState) GetFingerprint(trigger vos.Trigger) (vos.Fingerprint, bool) {
	fingerprint, ok := s.fingerprints[trigger]
	return fingerprint, ok
}

func (s *FingerprintState) StepName() string {
	return s.stepName
}

func (s *FingerprintState) Fingerprints() map[vos.Trigger]vos.Fingerprint {
	return s.fingerprints
}

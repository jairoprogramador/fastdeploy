package aggregates

import (
	"errors"
	"time"

	"github.com/jairoprogramador/fastdeploy/internal/domain/executionstate/vos"
)

// ScopeReceipt es un Agregado Raíz que actúa como un registro inmutable
// de que un paso específico se completó con éxito para un conjunto de fingerprints.
type ScopeReceipt struct {
	codeFingerprint        vos.Fingerprint
	environmentFingerprint vos.Fingerprint
	createdAt              time.Time
}

// NewScopeReceipt crea un nuevo y válido "recibo" de ejecución.
// Se asegura de que se proporcione al menos un fingerprint.
func NewScopeReceipt(
	codeFp vos.Fingerprint,
	envFp vos.Fingerprint,
) (*ScopeReceipt, error) {
	isCodeFpEmpty := (codeFp == vos.Fingerprint{})
	isEnvFpEmpty := (envFp == vos.Fingerprint{})

	if isCodeFpEmpty && isEnvFpEmpty {
		return nil, errors.New("un recibo de ejecución debe tener al menos un code fingerprint o un environment fingerprint")
	}

	return &ScopeReceipt{
		codeFingerprint:        codeFp,
		environmentFingerprint: envFp,
		createdAt:              time.Now().UTC(),
	}, nil
}

func NewScopeCodeReceipt(codeFp vos.Fingerprint) (*ScopeReceipt, error) {
	if (codeFp == vos.Fingerprint{}) {
		return nil, errors.New("el code fingerprint no puede estar vacío")
	}

	return &ScopeReceipt{
		codeFingerprint:        codeFp,
		environmentFingerprint: vos.Fingerprint{},
		createdAt:              time.Now().UTC(),
	}, nil
}

func NewScopeEnvironmentReceipt(envFp vos.Fingerprint) (*ScopeReceipt, error) {
	if (envFp == vos.Fingerprint{}) {
		return nil, errors.New("el environment fingerprint no puede estar vacío")
	}

	return &ScopeReceipt{
		codeFingerprint:        vos.Fingerprint{},
		environmentFingerprint: envFp,
		createdAt:              time.Now().UTC(),
	}, nil
}

func RehydrateScopeReceipt(
	codeFp vos.Fingerprint,
	envFp vos.Fingerprint,
	createdAt time.Time,
) (*ScopeReceipt, error) {
	isCodeFpEmpty := (codeFp == vos.Fingerprint{})
	isEnvFpEmpty := (envFp == vos.Fingerprint{})

	if isCodeFpEmpty && isEnvFpEmpty {
		return nil, errors.New("un recibo de ejecución debe tener al menos un code fingerprint o un environment fingerprint")
	}

	return &ScopeReceipt{
		codeFingerprint:        codeFp,
		environmentFingerprint: envFp,
		createdAt:              createdAt,
	}, nil
}

// CodeFingerprint devuelve el fingerprint del código asociado. Puede estar vacío.
func (r *ScopeReceipt) CodeFingerprint() vos.Fingerprint {
	return r.codeFingerprint
}

// EnvironmentFingerprint devuelve el fingerprint del ambiente asociado. Puede estar vacío.
func (r *ScopeReceipt) EnvironmentFingerprint() vos.Fingerprint {
	return r.environmentFingerprint
}

// CreatedAt devuelve la fecha y hora de creación del recibo.
func (r *ScopeReceipt) CreatedAt() time.Time {
	return r.createdAt
}

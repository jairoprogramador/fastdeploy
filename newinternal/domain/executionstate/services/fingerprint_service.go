package services

import (
	"context"

	"github.com/jairoprogramador/fastdeploy/newinternal/domain/executionstate/vos"
)

// FingerprintService define el contrato para un servicio de dominio que
// encapsula la lógica compleja de calcular los "fingerprints" (hashes)
// del estado del código y del ambiente.
type FingerprintService interface {
	// CalculateCodeFingerprint calcula un hash único para el estado actual del
	// código fuente de un proyecto, respetando las reglas de un archivo de ignore.
	CalculateCodeFingerprint(ctx context.Context, projectPath string, ignorePatterns []string) (vos.Fingerprint, error)

	// CalculateEnvironmentFingerprint calcula un hash único para el estado actual de
	// los archivos de configuración de infraestructura de un ambiente específico.
	CalculateEnvironmentFingerprint(ctx context.Context, environmentPath string) (vos.Fingerprint, error)
}

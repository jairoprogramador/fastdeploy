package vos

import "fmt"

// VerificationType define los tipos de hashes contra los que se puede verificar un paso.
type VerificationType int

const (
	VerificationTypeNone VerificationType = iota // No se verifica, siempre se ejecuta.
	VerificationTypeCode                         // Verificar contra el hash de código.
	VerificationTypeEnv                          // Verificar contra el hash de ambiente.
)

func (v VerificationType) String() string {
	switch v {
	case VerificationTypeCode:
		return "code"
	case VerificationTypeEnv:
		return "env"
	default:
		return "none"
	}
}

func VerificationTypeFromString(s string) (VerificationType, error) {
	switch s {
	case "code":
		return VerificationTypeCode, nil
	case "env":
		return VerificationTypeEnv, nil
	case "":
		return VerificationTypeNone, nil
	}
	return VerificationTypeNone, fmt.Errorf("tipo de verificación desconocido: '%s'", s)
}

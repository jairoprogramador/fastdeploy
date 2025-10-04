package vos

import "errors"

// Fingerprint representa el hash único del estado del código o de un ambiente.
// Es un Objeto de Valor inmutable.
type Fingerprint struct {
	value string
}

// NewFingerprint crea un nuevo y validado Objeto de Valor Fingerprint.
func NewFingerprint(value string) (Fingerprint, error) {
	if value == "" {
		return Fingerprint{}, errors.New("el valor del fingerprint no puede estar vacío")
	}
	return Fingerprint{value: value}, nil
}

// String devuelve la representación en cadena del fingerprint.
func (f Fingerprint) String() string {
	return f.value
}

// Equals comprueba la igualdad de valor con otro Fingerprint.
func (f Fingerprint) Equals(other Fingerprint) bool {
	return f.value == other.value
}

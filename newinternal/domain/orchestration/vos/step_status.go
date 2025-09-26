package vos

// StepStatus representa el estado de un StepExecution.
// Es un Objeto de Valor que utiliza un tipo enumerado.
type StepStatus int

const (
	// StepStatusPending indica que el paso aún no ha comenzado.
	StepStatusPending StepStatus = iota
	// StepStatusInProgress indica que el paso se está ejecutando actualmente.
	StepStatusInProgress
	// StepStatusSkipped indica que el paso fue omitido por el usuario.
	StepStatusSkipped
	// StepStatusSuccessful indica que el paso y todos sus comandos se completaron exitosamente.
	StepStatusSuccessful
	// StepStatusFailed indica que el paso falló porque uno de sus comandos falló.
	StepStatusFailed
)

// String devuelve la representación en cadena del estado.
func (s StepStatus) String() string {
	switch s {
	case StepStatusPending:
		return "Pendiente"
	case StepStatusInProgress:
		return "En Progreso"
	case StepStatusSkipped:
		return "Omitido"
	case StepStatusSuccessful:
		return "Exitoso"
	case StepStatusFailed:
		return "Fallido"
	default:
		return "Desconocido"
	}
}

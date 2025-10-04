package vos

// CommandStatus representa el estado de un CommandExecution.
// Es un Objeto de Valor que utiliza un tipo enumerado para mayor seguridad y claridad.
type CommandStatus int

const (
	// CommandStatusPending indica que el comando aún no se ha ejecutado.
	CommandStatusPending CommandStatus = iota
	// CommandStatusSuccessful indica que el comando se ejecutó y pasó todas sus validaciones de salida.
	CommandStatusSuccessful
	// CommandStatusFailed indica que el comando falló en su ejecución o en sus validaciones de salida.
	CommandStatusFailed
)

// String devuelve la representación en cadena del estado.
func (s CommandStatus) String() string {
	switch s {
	case CommandStatusPending:
		return "Pendiente"
	case CommandStatusSuccessful:
		return "Exitoso"
	case CommandStatusFailed:
		return "Fallido"
	default:
		return "Desconocido"
	}
}

func CommandStatusFromString(status string) CommandStatus {
	switch status {
	case "Pendiente":
		return CommandStatusPending
	case "Exitoso":
		return CommandStatusSuccessful
	case "Fallido":
		return CommandStatusFailed
	default:
		return CommandStatus(99)
	}
}
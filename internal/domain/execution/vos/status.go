package vos

type Status string

const (
	StatusPending    Status = "Pendiente"
	StatusInProgress Status = "En Progreso"
	StatusSkipped    Status = "Omitido"
	StatusCached     Status = "En Caché"
	StatusSuccessful Status = "Exitoso"
	StatusFailed     Status = "Fallido"
	StatusUnknown    Status = "Desconocido"
)

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "Pendiente"
	case StatusInProgress:
		return "En Progreso"
	case StatusSkipped:
		return "Omitido"
	case StatusCached:
		return "En Caché"
	case StatusSuccessful:
		return "Exitoso"
	case StatusFailed:
		return "Fallido"
	default:
		return "Desconocido"
	}
}
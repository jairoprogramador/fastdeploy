package vos

// StepStatus define el resultado de la ejecución de un paso.
type StepStatus string

const (
	// Success indica que el paso se ejecutó correctamente.
	Success StepStatus = "SUCCESS"
	// Failure indica que el paso falló durante la ejecución.
	Failure StepStatus = "FAILURE"
	// Cached indica que la ejecución del paso se omitió debido a la caché.
	Cached StepStatus = "CACHED"
)

// VariableSet representa un conjunto de variables, típicamente como un mapa de clave-valor.
type VariableSet map[string]string

// ExecutionResult encapsula el resultado de la ejecución de un único paso.
type ExecutionResult struct {
	// Status es el estado final del paso (éxito, fallo, etc.).
	Status StepStatus
	// Logs contiene la salida estándar y de error combinada del comando ejecutado.
	Logs string
	// OutputVars son las nuevas variables extraídas de la salida del comando.
	OutputVars VariableSet
	// Error contiene el error que causó el fallo, si lo hubo.
	Error error
}

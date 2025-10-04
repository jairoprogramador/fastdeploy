package dto

// StepDTO representa la estructura de los metadatos de un paso en un archivo step.yaml.
type StepDTO struct {
	Verifications []string `yaml:"verifications"`
}

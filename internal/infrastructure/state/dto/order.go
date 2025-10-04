package dto

// DTOs para la serialización del estado de la ejecución.

type CommandExecutionDTO struct {
	Name         string                          `yaml:"name"`
	Status       string                          `yaml:"status"`
	ResolvedCmd  string                          `yaml:"cmd"`
	ExecutionLog string                          `yaml:"log"`
	OutputVars   map[string]string                  `yaml:"outputs"`
}

type StepExecutionDTO struct {
	Name              string                `yaml:"name"`
	Status            string                `yaml:"status"`
	CommandExecutions []CommandExecutionDTO `yaml:"commands"`
}

type OrderDTO struct {
	ID                string                    `yaml:"id"`
	Status            string                    `yaml:"status"`
	TargetEnvironment string                    `yaml:"environment"`
	StepExecutions    []StepExecutionDTO        `yaml:"steps"`
	VariableMap       map[string]string   `yaml:"variables"`
}

package dto

import (
	deploymentvos "github.com/jairoprogramador/fastdeploy/newinternal/domain/deployment/vos"
	"github.com/jairoprogramador/fastdeploy/newinternal/domain/orchestration/vos"
)

// DTOs para la serialización del estado de la ejecución.

type CommandExecutionDTO struct {
	Name         string                          `yaml:"name"`
	Definition   deploymentvos.CommandDefinition `yaml:"definition"` // <-- AÑADIDO: El snapshot de la definición.
	Status       vos.CommandStatus               `yaml:"status"`
	ResolvedCmd  string                          `yaml:"resolved_cmd"`
	ExecutionLog string                          `yaml:"execution_log"`
	OutputVars   []vos.Variable                  `yaml:"output_vars"`
}

type StepExecutionDTO struct {
	Name              string                `yaml:"name"`
	Status            vos.StepStatus        `yaml:"status"`
	CommandExecutions []CommandExecutionDTO `yaml:"command_executions"`
}

type OrderDTO struct {
	ID                string                    `yaml:"id"`
	Status            vos.OrderStatus           `yaml:"status"`
	TargetEnvironment deploymentvos.Environment `yaml:"target_environment"`
	StepExecutions    []StepExecutionDTO        `yaml:"step_executions"`
	VariableMap       map[string]vos.Variable   `yaml:"variable_map"`
}
